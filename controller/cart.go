package controller

import (
	"context"
	"database/sql"
	"e-montir/api/handler"
	"e-montir/model"
	"e-montir/pkg/date"
	"e-montir/pkg/validator"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type cartCtx struct {
	CartModel model.Cart
	UserModel model.User
}

type Cart interface {
	SetCartAppointment(ctx context.Context, form *CartAppointmentRequest) error
	RemoveCartAppointment(ctx context.Context, userID string) error
	AddServiceToCart(ctx context.Context, serviceID int, cartID string) (*CartTotalItemAndPrice, error)
	RemoveServiceFromCart(ctx context.Context, serviceID int, cartID string) (*CartTotalItemAndPrice, error)
	CartDetail(ctx context.Context, userID string) (*CartDetail, error)
}

func NewCart(cartModel model.Cart, userModel model.User) Cart {
	return &cartCtx{
		CartModel: cartModel,
		UserModel: userModel,
	}
}

type (
	CartAppointmentRequest struct {
		UserID    string
		Date      string `json:"date"` // yyyy-mm-dd
		Time      string `json:"time"`
		BrandName string `json:"motorcycle_brand_name"`
	}

	CartAppointment struct {
		Date      string `json:"date"` // yyyy-mm-dd
		Time      string `json:"time"`
		BrandName string `json:"motorcycle_brand_name"`
	}

	AddOrRemoveServiceToCartRequest struct {
		ServiceID       int
		ServiceIDString string `json:"service_id"`
		CartID          int
		CartIDString    string
	}

	CartDetail struct {
		Location    UserCheckoutAddress `json:"location"`
		Appointment CartAppointment     `json:"appointment"`
		Items       []CartItems         `json:"items"`
		TotalPrice  float64             `json:"total_price"`
	}

	CartItems struct {
		CartID  int     `json:"id"`
		Title   string  `json:"title"`
		Price   float64 `json:"price"`
		Picture string  `json:"picture"`
	}

	CartTotalItemAndPrice struct {
		TotalPrice float64 `json:"total_price"`
		TotalItem  float64 `json:"total_item"`
	}
)

func (req *CartAppointmentRequest) ValidateCartAppointment() ([]handler.Fields, error) {
	var fields []handler.Fields
	var count int
	err := validator.ValidateID(req.UserID)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "user_id",
			Message: err.Error(),
		})
	}
	datePropareFormat, err := validator.ValidateDate(req.Date)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "date",
			Message: err.Error(),
		})
	}
	req.Date = datePropareFormat
	err = validator.ValidateTime(req.Time)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "time",
			Message: err.Error(),
		})
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

// nolint(gosec) // false positive
func (req *AddOrRemoveServiceToCartRequest) ValidateAddOrRemoveServiceToCart() ([]handler.Fields, error) {
	var fields []handler.Fields
	var count int
	serviceIDValid := true

	err := validator.ValidateServiceID(req.ServiceIDString)
	if err != nil {
		serviceIDValid = false
		count++
		fields = append(fields, handler.Fields{
			Name:    "service_id",
			Message: err.Error(),
		})
	}

	if serviceIDValid {
		service, serviceErr := strconv.Atoi(req.ServiceIDString)
		if serviceErr != nil {
			serviceIDValid = false
			count++
			fields = append(fields, handler.Fields{
				Name:    "service_id",
				Message: serviceErr.Error(),
			})
		}
		if service < 1 && serviceIDValid {
			count++
			fields = append(fields, handler.Fields{
				Name:    "service_id",
				Message: "service_id must be more than 0",
			})
		}
		req.ServiceID = service
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

// nolint(gosec) // false positive
func (c *cartCtx) SetCartAppointment(ctx context.Context, form *CartAppointmentRequest) error {
	isCartAvailable, data, err := c.CartModel.IsCartAvailable(ctx, form.UserID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when checking IsCartAvailable: %w", err)).Send()
		return err
	}

	if isCartAvailable {
		appDate, err := time.Parse(time.RFC3339, data.Date)
		if err != nil {
			log.Error().Err(fmt.Errorf("error when parsingDate: %w", err)).Send()
			return &handler.InternalServerError
		}
		// return nil because user is trying to create the same appointment as in the database
		if appDate.Local().Format(date.Format) == form.Date && data.Time == form.Time {
			return nil
		}

		// return error when user is trying to create new appointment while there is an appointment in database
		// user needs to remove the appointment first before creating a new one
		return &handler.CartAppointmentAvailable
	}

	err = c.CartModel.SetCartAppointment(ctx, &model.CartAppointment{
		UserID:    form.UserID,
		Date:      form.Date,
		Time:      form.Time,
		BrandName: form.BrandName,
	})

	if err != nil {
		log.Error().Err(fmt.Errorf("error when SetCartAppointment: %w", err)).Send()
		return err
	}

	return nil
}

func (c *cartCtx) RemoveCartAppointment(ctx context.Context, userID string) error {
	err := c.CartModel.RemoveCartAppointment(ctx, userID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when RemoveCartAppointment: %w", err)).Send()
		return err
	}

	return nil
}

func (c *cartCtx) AddServiceToCart(ctx context.Context, serviceID int, cartID string) (*CartTotalItemAndPrice, error) {
	isCartAvailable, _, err := c.CartModel.IsCartAvailable(ctx, cartID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when checking IsCartAvailable: %w", err)).Send()
		return nil, err
	}

	if !isCartAvailable {
		return nil, &handler.CartAppointmentNotAvailable
	}

	isServiceAvailable, err := c.CartModel.IsServiceAvailable(ctx, serviceID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when checking IsServiceAvailable: %w", err)).Send()
		return nil, err
	}

	if !isServiceAvailable {
		return nil, &handler.ServiceNotExists
	}

	res, err := c.CartModel.InsertServiceToCartItem(ctx, cartID, serviceID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when InsertServiceToCartItem: %w", err)).Send()
		return nil, err
	}
	return &CartTotalItemAndPrice{
		TotalPrice: res.TotalPrice,
		TotalItem:  res.TotalItem,
	}, nil
}

//nolint
func (c *cartCtx) RemoveServiceFromCart(ctx context.Context, serviceID int, cartID string) (*CartTotalItemAndPrice, error) {
	res, err := c.CartModel.RemoveServiceFromCartItem(ctx, serviceID, cartID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when RemoveServiceFromCartItem: %w", err)).Send()
		return nil, err
	}

	return &CartTotalItemAndPrice{
		TotalPrice: res.TotalPrice,
		TotalItem:  res.TotalItem,
	}, nil
}

func (c *cartCtx) CartDetail(ctx context.Context, userID string) (*CartDetail, error) {
	cartItem := make([]CartItems, 0)
	var totalPrice float64
	var cartDetail CartDetail
	res, err := c.CartModel.GetCartDetail(ctx, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(fmt.Errorf("error when GetCheckoutDetail: %w", err)).Send()
			return nil, err
		}

		cartDetail.Items = cartItem
		// when appointment is not been set
		return &cartDetail, nil
	}

	for _, v := range res.CartItem {
		cartItem = append(cartItem, CartItems{
			CartID:  v.CartID,
			Title:   v.Title,
			Price:   v.Price,
			Picture: v.Picture,
		})
		totalPrice += v.Price
	}

	loc, err := c.UserModel.GetUserCurrentLocation(ctx, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(fmt.Errorf("error when GetUserLocation: %w", err)).Send()
			return nil, err
		}

		// when location is not been set
		appointmentDate, dateParseErr := time.Parse(time.RFC3339, res.Appointment.Date)
		if dateParseErr != nil {
			log.Error().Err(fmt.Errorf("error when parsingDate: %w", err)).Send()
			return nil, &handler.InternalServerError
		}

		cartDetail.Appointment.Date = appointmentDate.Local().Format(date.Format)
		cartDetail.Appointment.Time = res.Appointment.Time
		cartDetail.Items = cartItem
		cartDetail.TotalPrice = totalPrice
		return &cartDetail, nil
	}

	// case when all data is available
	appointmentDate, err := time.Parse(time.RFC3339, res.Appointment.Date)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when parsingDate: %w", err)).Send()
		return nil, &handler.InternalServerError
	}

	cartDetail.Location.AddressID = loc.ID
	cartDetail.Location.Address = loc.Address
	cartDetail.Location.PhoneNum = loc.PhoneNumber
	cartDetail.Location.Label = loc.Label
	cartDetail.Location.Recipient = loc.RecipientName
	cartDetail.Appointment.Date = appointmentDate.Local().Format(date.Format)
	cartDetail.Appointment.Time = res.Appointment.Time
	cartDetail.Appointment.BrandName = res.Appointment.BrandName
	cartDetail.Items = cartItem
	cartDetail.TotalPrice = totalPrice
	return &cartDetail, nil
}
