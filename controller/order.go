package controller

import (
	"context"
	"database/sql"
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
	ListOfOrders(ctx context.Context, userID string) (*OrderListResponse, error)
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
		Title     string  `json:"title"`
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

	Mechanic struct {
		Name             string `json:"name"`
		PhoneNumber      string `json:"phone_number"`
		CompletedService int    `json:"completed_service"`
		Picture          string `json:"picture"`
	}

	OrderListData struct {
		ID              string           `json:"id"`
		UserID          string           `json:"user_id"`
		Description     string           `json:"description"`
		MotorCycleBrand string           `json:"motor_cycle_brand"`
		CreatedAt       string           `json:"created_at"`
		Appointment     OrderAppointment `json:"appointment"`
		Location        OrderLocation    `json:"location"`
		Items           []OrderItem      `json:"items"`
		Mechanic        Mechanic         `json:"mechanic"`
		TotalPrice      float64          `json:"total_price"`
		StatusOrder     string           `json:"status_order"`
		StatusDetail    string           `json:"status_detail"`
	}

	PlcaeOrderResponse struct {
		OrderID string `json:"order_id"`
	}

	OrderListResponse struct {
		Data []OrderListData `json:"data"`
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

	if totalPrice < 250000 {
		totalPrice = totalPrice + 15000
	}

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

func (c *orderCtx) ListOfOrders(ctx context.Context, userID string) (*OrderListResponse, error) {
	orderListData := make([]OrderListData, 0)
	orderLists, err := c.orderModel.ListOfOrders(ctx, userID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getListOfOrders : %w", err)).Send()
		return nil, err
	}

	for _, orderlist := range orderLists {
		var orderItems []OrderItem
		res, err := c.orderModel.ListOfOrderItems(ctx, orderlist.ID)
		if err != nil {
			log.Error().Err(fmt.Errorf("error when getListOfOrderItems : %w", err)).Send()
			return nil, err
		}

		for _, v := range res {
			tmp := OrderItem{
				ServiceID: v.ServiceID,
				Title:     v.Title,
				Price:     v.Price,
				Picture:   v.Picture,
			}
			orderItems = append(orderItems, tmp)
		}

		userLoc, err := c.orderModel.OrderLocation(ctx, orderlist.UserAddressID)
		if err != nil {
			log.Error().Err(fmt.Errorf("error when getOrderLocation : %w", err)).Send()
			return nil, err
		}

		appointmentDate, err := time.Parse(time.RFC3339, orderlist.Date)
		if err != nil {
			log.Error().Err(fmt.Errorf("error when parsingDate: %w", err)).Send()
			return nil, &handler.InternalServerError
		}

		appointment := OrderAppointment{
			Date: appointmentDate.Local().Format("2006-01-02"),
			Time: orderlist.TimeSlot,
		}

		mechanic, err := c.orderModel.GetOrderMechanic(ctx, int(orderlist.MechanicID.Int64))
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error().Err(fmt.Errorf("error when GetOrderMechanic: %w", err)).Send()
				return nil, &handler.InternalServerError
			}

			orderListData = append(orderListData, OrderListData{
				ID:              orderlist.ID,
				UserID:          orderlist.UserID,
				Description:     orderlist.Description.String,
				MotorCycleBrand: orderlist.MotorCycleBrand,
				Appointment:     appointment,
				Location: OrderLocation{
					AddressID: userLoc.ID,
					Address:   userLoc.Address,
					Label:     userLoc.Label,
					Recipient: userLoc.RecipientName,
					PhoneNum:  userLoc.PhoneNumber,
				},
				CreatedAt:    orderlist.CreatedAt.Format(time.RFC3339),
				Mechanic:     Mechanic{},
				Items:        orderItems,
				TotalPrice:   orderlist.TotalPrice,
				StatusOrder:  orderlist.OrderStatus.String,
				StatusDetail: orderlist.OrderStatus.String,
			})
		} else {
			orderListData = append(orderListData, OrderListData{
				ID:              orderlist.ID,
				UserID:          orderlist.UserID,
				Description:     orderlist.Description.String,
				MotorCycleBrand: orderlist.MotorCycleBrand,
				Appointment:     appointment,
				Location: OrderLocation{
					AddressID: userLoc.ID,
					Address:   userLoc.Address,
					Label:     userLoc.Label,
					Recipient: userLoc.RecipientName,
					PhoneNum:  userLoc.PhoneNumber,
				},
				Mechanic: Mechanic{
					Name:             mechanic.Name,
					PhoneNumber:      mechanic.PhoneNumber,
					CompletedService: mechanic.CompletedService,
					Picture:          mechanic.Picture.String,
				},
				CreatedAt:    orderlist.CreatedAt.Format(time.RFC3339),
				Items:        orderItems,
				TotalPrice:   orderlist.TotalPrice,
				StatusOrder:  orderlist.OrderStatus.String,
				StatusDetail: orderlist.OrderStatus.String,
			})
		}
	}

	return &OrderListResponse{
		Data: orderListData,
	}, nil
}
