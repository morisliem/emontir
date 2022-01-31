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
	orderModel  model.Order
	cartModel   model.Cart
	userModel   model.User
	reviewModel model.Review
}

type Order interface {
	PlaceOrder(ctx context.Context, userID, orderID, invoiceID string) (*PlcaeOrderResponse, error)
	PaymentReceived(ctx context.Context, orderID, transactionStatus string) error
	ListOfOrders(ctx context.Context, userID string) (*OrderListResponse, error)
	UpdateOrderStatus(ctx context.Context, form *UpdateOrderRequest) error
	OrderDetail(ctx context.Context, orderID string) (*OrderDetailResponse, error)
}

func NewOrder(orderModel model.Order, cartModel model.Cart, userModel model.User, reviewModel model.Review) Order {
	return &orderCtx{
		orderModel:  orderModel,
		cartModel:   cartModel,
		userModel:   userModel,
		reviewModel: reviewModel,
	}
}

type (
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
		InvoiceID       string           `json:"invoice_id"`
		IsReviewed      bool             `json:"is_reviewed"`
	}

	PlcaeOrderResponse struct {
		OrderID string `json:"order_id"`
	}

	OrderListResponse struct {
		Data []OrderListData `json:"data"`
	}

	OrderDetailResponse struct {
		Data OrderListData `json:"data"`
	}

	UpdateOrderRequest struct {
		ID     string `json:"invoice_id"`
		Status string `json:"status"`
	}
)

func (c *orderCtx) PlaceOrder(ctx context.Context, userID, orderID, invoiceID string) (*PlcaeOrderResponse, error) {
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
		InvoiceID:       invoiceID,
	})

	if err != nil {
		log.Error().Err(fmt.Errorf("error when setOrder : %w", err)).Send()
		return nil, err
	}

	return &PlcaeOrderResponse{
		OrderID: orderID,
	}, nil
}

func (c *orderCtx) PaymentReceived(ctx context.Context, orderID, transactionStatus string) error {
	// var notif fcm.NotifFcm
	// userID, err := c.userModel.GetUserIDByOrderID(ctx, orderID)
	// if err != nil {
	// 	log.Error().Err(fmt.Errorf("error when GetUserIDByOrderID : %w", err)).Send()
	// 	return err
	// }

	// fcmKey, err := c.userModel.GetFCMKey(ctx, userID)
	// if err != nil {
	// 	log.Error().Err(fmt.Errorf("error when GetFCMKey : %w", err)).Send()
	// 	return err
	// }

	// if transactionStatus != "PAID" {
	// 	notif.To = fcmKey
	// 	notif.Title = "payment failed"
	// 	notif.Body = "payment failed, please try again"
	// 	notif.Redirect = fmt.Sprintf("%s/orders/{%s}", os.Getenv("BASE_URL"), orderID)
	// 	fcm.SendNotification(ctx, notif)
	// 	return fmt.Errorf("payment failed")
	// }

	err := c.orderModel.UpdateOrderStatus(ctx, orderID, "On process", "")
	if err != nil {
		log.Error().Err(fmt.Errorf("error when UpdateOrderStatus : %w", err)).Send()
		return err
	}

	// notif.To = fcmKey
	// notif.Title = "payment success"
	// notif.Body = "preparing your order"
	// notif.Redirect = fmt.Sprintf("%s/orders/{%s}", os.Getenv("BASE_URL"), orderID)
	// fcm.SendNotification(ctx, notif)

	err = c.orderModel.AssignMechanic(ctx, orderID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when assignMechanic : %w", err)).Send()
		return err
	}

	return nil
}

func (c *orderCtx) ListOfOrders(ctx context.Context, userID string) (*OrderListResponse, error) {
	orderListData := make([]OrderListData, 0)
	orderLists, err := c.orderModel.ListOfOrders(ctx, userID)
	isReviewed := false
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

		_, err = c.reviewModel.GetReviewByOrderID(ctx, orderlist.ID)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error().Err(fmt.Errorf("error when GetReviewByOrderID: %w", err)).Send()
				return nil, &handler.InternalServerError
			}
		}

		if err == nil {
			isReviewed = true
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
				StatusDetail: orderlist.OrderDetail.String,
				InvoiceID:    orderlist.InvoiceID,
				IsReviewed:   isReviewed,
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
				StatusDetail: orderlist.OrderDetail.String,
				InvoiceID:    orderlist.InvoiceID,
				IsReviewed:   isReviewed,
			})
		}
	}

	return &OrderListResponse{
		Data: orderListData,
	}, nil
}

func (c *orderCtx) UpdateOrderStatus(ctx context.Context, form *UpdateOrderRequest) error {
	// userID, orderID, err := c.userModel.GetUserIDNOrderIDByInvoiceID(ctx, form.ID)
	// if err != nil {
	// 	log.Error().Err(fmt.Errorf("error when GetUserIDByOrderID : %w", err)).Send()
	// 	return err
	// }

	// fcmKey, err := c.userModel.GetFCMKey(ctx, userID)
	// if err != nil {
	// 	log.Error().Err(fmt.Errorf("error when GetFCMKey : %w", err)).Send()
	// 	return err
	// }

	err := c.orderModel.UpdateOrderStatus(ctx, "", form.Status, form.ID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when UpdateOrderStatus : %w", err)).Send()
		return err
	}

	if form.Status == "on the way" {
		// fcm.SendNotification(ctx, fcm.NotifFcm{
		// 	To:       fcmKey,
		// 	Redirect: fmt.Sprintf("%s/orders/{%s}", os.Getenv("BASE_URL"), orderID),
		// 	Title:    "mechanic is on the way",
		// 	Body:     "mechanic is on the way to your place. please wait",
		// })
		return nil
	}

	if form.Status == "done" {
		err = c.orderModel.OrderCompleted(ctx, form.ID)
		if err != nil {
			log.Error().Err(fmt.Errorf("error when execute OrderCompleted : %w", err)).Send()
			return err
		}
		// fcm.SendNotification(ctx, fcm.NotifFcm{
		// 	To:       fcmKey,
		// 	Redirect: fmt.Sprintf("%s/orders/{%s}", os.Getenv("BASE_URL"), orderID),
		// 	Title:    "service done",
		// 	Body:     "service done, looking forward to your next order",
		// })
		return nil
	}
	return nil
}

func (c *orderCtx) OrderDetail(ctx context.Context, orderID string) (*OrderDetailResponse, error) {
	var orderDetailResponse OrderDetailResponse
	var orderItems []OrderItem
	isReviewed := false

	orderDetail, err := c.orderModel.GetOrderByOrderID(ctx, orderID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getOrderByOrderID : %w", err)).Send()
		return nil, err
	}

	listOfOrderItem, err := c.orderModel.ListOfOrderItems(ctx, orderID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getListOfOrderItems : %w", err)).Send()
		return nil, err
	}

	for _, v := range listOfOrderItem {
		tmp := OrderItem{
			ServiceID: v.ServiceID,
			Title:     v.Title,
			Price:     v.Price,
			Picture:   v.Picture,
		}
		orderItems = append(orderItems, tmp)
	}

	userLoc, err := c.orderModel.OrderLocation(ctx, orderDetail.UserAddressID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when getOrderLocation : %w", err)).Send()
		return nil, err
	}

	appointmentDate, err := time.Parse(time.RFC3339, orderDetail.Date)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when parsingDate: %w", err)).Send()
		return nil, &handler.InternalServerError
	}

	appointment := OrderAppointment{
		Date: appointmentDate.Local().Format("2006-01-02"),
		Time: orderDetail.TimeSlot,
	}

	_, err = c.reviewModel.GetReviewByOrderID(ctx, orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(fmt.Errorf("error when GetReviewByOrderID: %w", err)).Send()
			return nil, &handler.InternalServerError
		}
	}

	if err == nil {
		isReviewed = true
	}

	mechanic, err := c.orderModel.GetOrderMechanic(ctx, int(orderDetail.MechanicID.Int64))
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(fmt.Errorf("error when GetOrderMechanic: %w", err)).Send()
			return nil, &handler.InternalServerError
		}

		orderDetailResponse.Data.ID = orderID
		orderDetailResponse.Data.UserID = orderDetail.UserID
		orderDetailResponse.Data.Description = orderDetail.Description.String
		orderDetailResponse.Data.MotorCycleBrand = orderDetail.MotorCycleBrand
		orderDetailResponse.Data.Appointment = appointment
		orderDetailResponse.Data.Location = OrderLocation{
			AddressID: userLoc.ID,
			Address:   userLoc.Address,
			Label:     userLoc.Label,
			Recipient: userLoc.RecipientName,
			PhoneNum:  userLoc.PhoneNumber,
		}
		orderDetailResponse.Data.CreatedAt = orderDetail.CreatedAt.Format(time.RFC3339)
		orderDetailResponse.Data.Mechanic = Mechanic{}
		orderDetailResponse.Data.Items = orderItems
		orderDetailResponse.Data.TotalPrice = orderDetail.TotalPrice
		orderDetailResponse.Data.StatusOrder = orderDetail.OrderStatus.String
		orderDetailResponse.Data.StatusDetail = orderDetail.OrderDetail.String
		orderDetailResponse.Data.InvoiceID = orderDetail.InvoiceID
		orderDetailResponse.Data.IsReviewed = isReviewed
	} else {

		orderDetailResponse.Data.ID = orderID
		orderDetailResponse.Data.UserID = orderDetail.UserID
		orderDetailResponse.Data.Description = orderDetail.Description.String
		orderDetailResponse.Data.MotorCycleBrand = orderDetail.MotorCycleBrand
		orderDetailResponse.Data.Appointment = appointment
		orderDetailResponse.Data.Location = OrderLocation{
			AddressID: userLoc.ID,
			Address:   userLoc.Address,
			Label:     userLoc.Label,
			Recipient: userLoc.RecipientName,
			PhoneNum:  userLoc.PhoneNumber,
		}
		orderDetailResponse.Data.CreatedAt = orderDetail.CreatedAt.Format(time.RFC3339)
		orderDetailResponse.Data.Mechanic = Mechanic{
			Name:             mechanic.Name,
			PhoneNumber:      mechanic.PhoneNumber,
			CompletedService: mechanic.CompletedService,
			Picture:          mechanic.Picture.String,
		}
		orderDetailResponse.Data.Items = orderItems
		orderDetailResponse.Data.TotalPrice = orderDetail.TotalPrice
		orderDetailResponse.Data.StatusOrder = orderDetail.OrderStatus.String
		orderDetailResponse.Data.StatusDetail = orderDetail.OrderDetail.String
		orderDetailResponse.Data.InvoiceID = orderDetail.InvoiceID
		orderDetailResponse.Data.IsReviewed = isReviewed
	}

	return &orderDetailResponse, nil
}
