package model

import (
	"context"
	"database/sql"
	"e-montir/api/handler"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	OrderBaseModel struct {
		ID              string         `db:"id"`
		UserID          string         `db:"user_id"`
		UserAddressID   string         `db:"user_address_id"`
		Description     sql.NullString `db:"description"`
		TotalPrice      float64        `db:"total_price"`
		CreatedAt       time.Time      `db:"created_at"`
		OrderStatus     sql.NullString `db:"status_order"` // waiting for payment, on process, on the way, arrived, done
		OrderDetail     sql.NullString `db:"status_detail"`
		MotorCycleBrand string         `db:"motor_cycle_brand_name"`
		TimeSlot        string         `db:"time_slot"`
		Date            string         `db:"date"`
		MechanicID      sql.NullInt64  `db:"mechanic_id"`
	}

	OrderItem struct {
		ServiceID int     `db:"id"`
		Title     string  `db:"title"`
		Price     float64 `db:"price"`
		Picture   string  `db:"picture"`
	}

	OrderLocation struct {
		ID            string    `db:"id"`
		Label         string    `db:"label"`
		Address       string    `db:"address"`
		AddressDetail string    `db:"address_detail"`
		PhoneNumber   string    `db:"phone_num"`
		RecipientName string    `db:"recipient_name"`
		Latitude      string    `db:"latitude"`
		Longitude     string    `db:"longitude"`
		CreatedAt     time.Time `db:"created_at"`
	}

	OrderMechanic struct {
		ID               string         `db:"id"`
		Name             string         `db:"name"`
		PhoneNumber      string         `db:"phone_number"`
		CompletedService int            `db:"completed_service"`
		Picture          sql.NullString `db:"picture"`
		Status           bool           `db:"status"`
	}
)

type Order interface {
	SetOrder(ctx context.Context, userID string, param *OrderBaseModel) error
	AssignMechanic(ctx context.Context, orderID string) error
	CheckOrder(ctx context.Context, orderID string) (*OrderBaseModel, error)
	ListOfOrders(ctx context.Context, userID string) ([]OrderBaseModel, error)
	ListOfOrderItems(ctx context.Context, orderID string) ([]OrderItem, error)
	OrderLocation(ctx context.Context, userAddressID string) (*OrderLocation, error)
	GetOrderMechanic(ctx context.Context, mechanicID int) (*OrderMechanic, error)
}

type order struct {
	db      *sqlx.DB
	queries map[string]*sqlx.Stmt
}

func NewOrder(db *sqlx.DB) Order {
	order := new(order)
	order.db = db
	order.queries = make(map[string]*sqlx.Stmt, len(orderQueries))
	for k, v := range orderQueries {
		stmt, err := db.Preparex(v)
		if err != nil {
			log.Fatal().Msg("error : " + err.Error() + "\norder : " + v)
		}
		order.queries[k] = stmt
	}
	return order
}

var (
	setOrder        = "setOrder"
	setOrderField1  = `("id", "user_id", "user_address_id", "date", "time_slot", "created_at", `
	setOrderFields2 = `"total_price", "motor_cycle_brand_name", "status_order")`
	setOrderFields  = setOrderField1 + setOrderFields2
	setOrderSQL     = `INSERT INTO "orders" ` + setOrderFields + ` VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	getServiceIDSQL = `SELECT "service_id" FROM "cart_items" WHERE "cart_id" = $1`

	insertOrderItemSQL = `INSERT INTO "order_items" ("service_id", "order_id") VALUES ($1,$2)`

	updateEmployeeNum          = "updateEmployeeNum"
	updateEmployeeNumCondition = `WHERE "date" = $1 AND "time" = $2 AND "employee_num" > 0`
	updateEmployeeNumSQL       = `UPDATE "time_slots" SET "employee_num"="employee_num" - 1 ` + updateEmployeeNumCondition

	checkEmployeeAvailability    = "employeeAvailability"
	checkEmployeeAvailabilitySQL = `SELECT "employee_num" FROM "time_slots" WHERE "date"=$1 AND "time"=$2`

	removeOrder    = "removeOrder"
	removeOrderSQL = `DELETE FROM "orders" WHERE id = $1`

	assignMechanic    = "assignMontirToOrder"
	assignMechanicSQL = `UPDATE "orders" SET "mechanic_id" = $2, "status_order" = $3 WHERE "id" = $1`

	getMechanicIDs    = "getMechanic"
	getMechanicIDsSQL = `SELECT "id" FROM "mechanics" ORDER BY "is_available" DESC`

	updateMechanicAvailability    = "updateMechanicAvailability"
	updateMechanicAvailabilitySQL = `UPDATE "mechanics" SET "is_available" = $2 WHERE "id" = $1`

	getOrderItems     = "OrderItems"
	getOrderItemsJoin = `LEFT OUTER JOIN "services" ON order_items.service_id = services.id WHERE "order_id" = $1 `
	getOrderItemsSQL  = `SELECT services.id AS "id", "title", "price", "picture" FROM "order_items" ` + getOrderItemsJoin

	getOrderLocation          = "getOrderLocation"
	getOrderLocationFields    = `"id","label","address","address_detail","phone_num","recipient_name","latitude","longitude"`
	getOrderLocationCondition = `WHERE "id" = $1`
	getOrderLocationSQL       = `SELECT ` + getOrderLocationFields + ` FROM "user_addresses" ` + getOrderLocationCondition

	getOrderListByID    = "getOrder"
	getOrderListField1  = `"id", "description", "total_price", "user_address_id", "created_at", "status_detail", `
	getOrderListField2  = `"status_order", "user_id", "motor_cycle_brand_name", "time_slot", "date", "mechanic_id"`
	getOrderListField   = getOrderListField1 + getOrderListField2
	getOrderListByIDSQL = `SELECT ` + getOrderListField + `FROM "orders" WHERE "id" = $1`

	getOrderListByUserID    = "getOrder"
	getOrderListByUserIDSQL = `SELECT ` + getOrderListField + `FROM "orders" WHERE "user_id" = $1`

	getMechanic       = "getMechanic"
	getMechanicFields = `"id", "name", "phone_number", "completed_service", "picture"`
	getMechanicSQL    = `SELECT ` + getMechanicFields + ` FROM "mechanics" WHERE "id" = $1`

	orderQueries = map[string]string{
		setOrder:                   setOrderSQL,
		updateEmployeeNum:          updateEmployeeNumSQL,
		checkEmployeeAvailability:  checkEmployeeAvailabilitySQL,
		removeOrder:                removeOrderSQL,
		assignMechanic:             assignMechanicSQL,
		getMechanicIDs:             getMechanicIDsSQL,
		updateMechanicAvailability: updateMechanicAvailabilitySQL,
		getOrderListByID:           getOrderListByIDSQL,
		getOrderItems:              getOrderItemsSQL,
		getOrderLocation:           getOrderLocationSQL,
		getMechanic:                getMechanicSQL,
		getOrderListByUserID:       getOrderListByUserIDSQL,
	}
)

func (c *order) SetOrder(ctx context.Context, userID string, param *OrderBaseModel) error {
	var serviceIDs []string
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if rollback := tx.Rollback(); rollback == nil {
			log.Info().Msg("rolling back changes")
		}
	}()

	var employeeNum int
	err = tx.QueryRowContext(ctx, checkEmployeeAvailabilitySQL, param.Date, param.TimeSlot).Scan(&employeeNum)
	if err != nil {
		return err
	}

	if employeeNum <= 0 {
		return &handler.NoEmployeeError
	}

	rows, err := tx.QueryContext(ctx, getServiceIDSQL, userID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var serviceID string
		err = rows.Scan(&serviceID)
		if err != nil {
			return err
		}
		serviceIDs = append(serviceIDs, serviceID)
	}

	// nolint(gosec) // false positive
	_, err = tx.ExecContext(ctx, setOrderSQL, param.ID, userID, param.UserAddressID, param.Date, param.TimeSlot, param.CreatedAt, param.TotalPrice, param.MotorCycleBrand, "waiting for payment")
	if err != nil {
		return err
	}

	for _, serviceID := range serviceIDs {
		_, insertErr := tx.ExecContext(ctx, insertOrderItemSQL, serviceID, param.ID)
		if insertErr != nil {
			return insertErr
		}
	}

	_, err = tx.ExecContext(ctx, updateEmployeeNumSQL, param.Date, param.TimeSlot)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, removeCartAppointmentSQL, param.UserID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *order) AssignMechanic(ctx context.Context, orderID string) error {
	var mechanicID int
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if rollback := tx.Rollback(); rollback == nil {
			log.Info().Msg("rolling back changes")
		}
	}()

	err = tx.QueryRowContext(ctx, getMechanicIDsSQL).Scan(&mechanicID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, assignMechanicSQL, orderID, mechanicID, "on process")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, updateMechanicAvailabilitySQL, mechanicID, false)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *order) CheckOrder(ctx context.Context, orderID string) (*OrderBaseModel, error) {
	var order OrderBaseModel
	fmt.Println(orderID)
	err := c.queries[getOrderListByID].GetContext(ctx, &order, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &handler.OrderNotExists
		}
		return nil, err
	}

	if order.OrderStatus.String != "waiting for payment" {
		return nil, &handler.OrderHasBeenPaid
	}

	return &order, nil
}

func (c *order) ListOfOrders(ctx context.Context, userID string) ([]OrderBaseModel, error) {
	var result []OrderBaseModel
	err := c.queries[getOrderListByUserID].SelectContext(ctx, &result, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *order) ListOfOrderItems(ctx context.Context, orderID string) ([]OrderItem, error) {
	var result []OrderItem
	err := c.queries[getOrderItems].SelectContext(ctx, &result, orderID)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (c *order) OrderLocation(ctx context.Context, userAddressID string) (*OrderLocation, error) {
	var result OrderLocation
	err := c.queries[getOrderLocation].GetContext(ctx, &result, userAddressID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *order) GetOrderMechanic(ctx context.Context, mechanicID int) (*OrderMechanic, error) {
	var result OrderMechanic
	err := c.queries[getMechanic].GetContext(ctx, &result, mechanicID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
