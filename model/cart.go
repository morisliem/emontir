package model

import (
	"context"
	"database/sql"
	"e-montir/api/handler"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	CartBaseModel struct {
		Appointment CartAppointment
		CartItem    []CartItemBaseModel
	}

	CartItemBaseModel struct {
		CartID  int     `db:"id"`
		Title   string  `db:"title"`
		Price   float64 `db:"price"`
		Picture string  `db:"picture"`
	}

	CartAppointment struct {
		UserID string `db:"user_id"`
		Date   string `db:"date"`
		Time   string `db:"time_slot"`
	}

	CartItemAndPrice struct {
		TotalItem  float64 `db:"total_item"`
		TotalPrice float64 `db:"total_price"`
	}
)

type Cart interface {
	SetCartAppointment(ctx context.Context, param *CartAppointment) error
	RemoveCartAppointment(ctx context.Context, cartID string) error
	InsertServiceToCartItem(ctx context.Context, cartID string, serviceID int) (*CartItemAndPrice, error)
	RemoveServiceFromCartItem(ctx context.Context, serviceID int, cartID string) (*CartItemAndPrice, error)
	GetCartDetail(ctx context.Context, uid string) (*CartBaseModel, error)
	IsCartAvailable(ctx context.Context, cartID string) (bool, *CartAppointment, error)
	IsServiceAvailable(ctx context.Context, serviceID int) (bool, error)
}

type cart struct {
	db      *sqlx.DB
	queries map[string]*sqlx.Stmt
}

func NewCart(db *sqlx.DB) Cart {
	cart := new(cart)
	cart.db = db
	cart.queries = make(map[string]*sqlx.Stmt, len(CartQueries))
	for k, v := range CartQueries {
		stmt, err := db.Preparex(v)
		if err != nil {
			log.Fatal().Msg("error : " + err.Error() + "\nCart : " + v)
		}
		cart.queries[k] = stmt
	}
	return cart
}

var (
	getTotalPriceAndTotalItem       = "totalPriceAndTotalItem"
	getTotalPriceAndTotalItemSelect = `SELECT COUNT(*) as total_item, SUM(price) as total_price FROM`
	getTotalPriceAndTotalItemJoin   = `LEFT OUTER JOIN "services" ON cart_items.service_id=services.id WHERE "cart_id"=$1`
	getTotalPriceAndTotalItemSQL    = getTotalPriceAndTotalItemSelect + " cart_items " + getTotalPriceAndTotalItemJoin

	setAppointment    = "setAppointment"
	setAppointmentSQL = `INSERT INTO "carts" ("id", "user_id", "date", "time_slot") VALUES ($1,$2,$3,$4)`

	getAppointment    = "getAppointment"
	getAppointmentSQL = `SELECT "date", "time_slot" FROM "carts" WHERE "user_id" = $1`

	removeCartAppointment    = "removeCartAppointment"
	removeCartAppointmentSQL = `DELETE FROM "carts" WHERE "id" = $1`

	insertServiceToCartItem    = "addService"
	insertServiceToCartItemSQL = `INSERT INTO "cart_items" ("cart_id", "service_id") VALUES ($1,$2)`

	removeServiceFromCartItem    = "removeService"
	removeServiceFromCartItemSQL = `DELETE FROM "cart_items" WHERE "service_id" = $1 AND "cart_id" = $2`

	checkServiceAvailability    = "serviceAvailability"
	checkServiceAvailabilitySQL = `SELECT "title" FROM "services" WHERE "id" = $1`

	getCartItems     = "CartItems"
	getCartItemsJoin = `LEFT OUTER JOIN "services" ON cart_items.service_id = services.id WHERE "cart_id" = $1 `
	getCartItemsSQL  = `SELECT services.id AS "id", "title", "price", "picture" FROM "cart_items" ` + getCartItemsJoin
	CartQueries      = map[string]string{
		setAppointment:            setAppointmentSQL,
		getAppointment:            getAppointmentSQL,
		removeCartAppointment:     removeCartAppointmentSQL,
		insertServiceToCartItem:   insertServiceToCartItemSQL,
		removeServiceFromCartItem: removeServiceFromCartItemSQL,
		getCartItems:              getCartItemsSQL,
		getTotalPriceAndTotalItem: getTotalPriceAndTotalItemSQL,
		checkServiceAvailability:  checkServiceAvailabilitySQL,
	}
)

func (c *cart) SetCartAppointment(ctx context.Context, param *CartAppointment) error {
	_, err := c.queries[setAppointment].ExecContext(ctx, param.UserID, param.UserID, param.Date, param.Time)
	if err != nil {
		return err
	}
	return nil
}

func (c *cart) RemoveCartAppointment(ctx context.Context, cartID string) error {
	row, err := c.queries[removeCartAppointment].ExecContext(ctx, cartID)
	if err != nil {
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil || rowAffected != 1 {
		return &handler.CartAppointmentNotAvailable
	}

	return nil
}

func (c *cart) InsertServiceToCartItem(ctx context.Context, cartID string, serviceID int) (*CartItemAndPrice, error) {
	_, err := c.queries[insertServiceToCartItem].ExecContext(ctx, cartID, serviceID)
	if err != nil {
		return nil, err
	}

	res, err := c.getTotalItemAndTotalPriceFromCart(ctx, cartID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *cart) RemoveServiceFromCartItem(ctx context.Context, serviceID int, cartID string) (*CartItemAndPrice, error) {
	row, err := c.queries[removeServiceFromCartItem].ExecContext(ctx, serviceID, cartID)
	if err != nil {
		return nil, err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil || rowAffected < 1 {
		return nil, &handler.ServiceNotExists
	}

	res, err := c.getTotalItemAndTotalPriceFromCart(ctx, cartID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *cart) GetCartDetail(ctx context.Context, uid string) (*CartBaseModel, error) {
	var appointment CartAppointment
	var cartItems []CartItemBaseModel
	err := c.queries[getAppointment].QueryRowContext(ctx, uid).Scan(&appointment.Date, &appointment.Time)
	if err != nil {
		return nil, err
	}

	err = c.queries[getCartItems].SelectContext(ctx, &cartItems, uid)
	if err != nil {
		return nil, err
	}

	return &CartBaseModel{
		Appointment: appointment,
		CartItem:    cartItems,
	}, nil
}

func (c *cart) getTotalItemAndTotalPriceFromCart(ctx context.Context, cartID string) (*CartItemAndPrice, error) {
	var totalPrice, totalItem *sql.NullFloat64
	err := c.queries[getTotalPriceAndTotalItem].QueryRowContext(ctx, cartID).Scan(&totalItem, &totalPrice)
	if err != nil {
		return nil, err
	}

	if totalPrice == nil {
		return &CartItemAndPrice{
			TotalPrice: 0,
			TotalItem:  totalItem.Float64,
		}, nil
	}

	return &CartItemAndPrice{
		TotalPrice: totalPrice.Float64,
		TotalItem:  totalItem.Float64,
	}, nil
}

func (c *cart) IsCartAvailable(ctx context.Context, cartID string) (bool, *CartAppointment, error) {
	var appointment CartAppointment
	err := c.queries[getAppointment].QueryRowContext(ctx, cartID).Scan(&appointment.Date, &appointment.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, &appointment, nil
}

func (c *cart) IsServiceAvailable(ctx context.Context, serviceID int) (bool, error) {
	var title string
	err := c.queries[checkServiceAvailability].QueryRowContext(ctx, serviceID).Scan(&title)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
