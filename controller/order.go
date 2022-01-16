package controller

import (
	"context"
	"e-montir/api/handler"
	"e-montir/model"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type orderCtx struct {
	orderModel model.Order
	cartModel  model.Cart
	userModel  model.User
}

type Order interface {
	PlaceOrder(ctx context.Context, userID, orderID string) (*PlcaeOrderResponse, error)
	PaymentReceived(ctx context.Context, orderID string) error
}

func NewOrder(orderModel model.Order, cartModel model.Cart, userModel model.User) Order {
	return &orderCtx{
		orderModel: orderModel,
		cartModel:  cartModel,
		userModel:  userModel,
	}
}

type (
	OrderResponse struct {
		ID              string    `json:"id"`
		UserID          string    `json:"user_id"`
		UserAddressID   string    `json:"user_address_id"`
		Description     string    `json:"description"`
		TotalPrice      float64   `json:"total_price"`
		CreatedAt       time.Time `json:"created_at"`
		Status          string    `json:"status"` // waiting for payment, paid, completed
		MotorCycleBrand string    `json:"motor_cycle_brand_name"`
		TimeSlot        string    `json:"time_slot:"`
		Date            string    `json:"date"`
	}

	OrderAppointment struct {
		Date string `json:"date"`
		Time string `json:"time"`
	}

	OrderItem struct {
		ServiceID int     `json:"id"`
		Title     string  `json:"string"`
		Price     float64 `json:"price"`
		Picture   string  `json:"picture"`
	}

	OrderLocation struct {
		AddressID string `json:"id"`
		Address   string `json:"address"`
		Label     string `json:"label"`
		Recipient string `json:"recipient"`
		PhoneNum  string `json:"phone_number"`
	}

	OrderListResponse struct {
		ID              string           `json:"id"`
		UserID          string           `json:"user_id"`
		Description     string           `json:"description"`
		MotorCycleBrand string           `json:"motor_cycle_brand"`
		Appointment     OrderAppointment `json:"appointment"`
		Location        OrderLocation    `json:"location"`
		Items           []OrderItem      `json:"items"`
		TotalPrice      float64          `json:"total_price"`
		Status          string           `json:"status"`
	}

	UpdateOrderRequest struct {
		ID string `json:"id"`
	}

	PlcaeOrderResponse struct {
		OrderID string `json:"order_id"`
	}
)

func (c *orderCtx) PlaceOrder(ctx context.Context, userID, orderID string) (*PlcaeOrderResponse, error) {
	isCartAvailable, _, err := c.cartModel.IsCartAvailable(ctx, userID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when checking IsCartAvailable: %w", err)).Send()
		return nil, err
	}

	if !isCartAvailable {
		return nil, &handler.CartAppointmentNotAvailable
	}

	totalPrice := 0
	res, err := c.cartModel.GetCartDetail(ctx, userID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getCartDetail : %w", err)).Send()
		return nil, err
	}

	for _, v := range res.CartItem {
		totalPrice += int(v.Price)
	}

	userLoc, err := c.userModel.GetUserCurrentLocation(ctx, userID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getUserCurrentLocation : %w", err)).Send()
		return nil, err
	}

	fmt.Println(res)

	err = c.orderModel.SetOrder(ctx, userID, &model.OrderBaseModel{
		ID:              orderID,
		UserID:          userID,
		UserAddressID:   userLoc.ID,
		TimeSlot:        res.Appointment.Time,
		Date:            res.Appointment.Date,
		MotorCycleBrand: res.Appointment.BrandName,
		TotalPrice:      float64(totalPrice),
		CreatedAt:       time.Now(),
	})

	if err != nil {
		log.Error().Err(fmt.Errorf("error when setOrder : %w", err)).Send()
		return nil, err
	}

	return &PlcaeOrderResponse{
		OrderID: orderID,
	}, nil
}

func (c *orderCtx) PaymentReceived(ctx context.Context, orderID string) error {
	err := c.orderModel.AssignMechanic(ctx, orderID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when assignMechanic : %w", err)).Send()
		return err
	}
	return nil
}
