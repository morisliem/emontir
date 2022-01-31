package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"e-montir/pkg/uuid"
	"e-montir/pkg/validator"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type OrderHandler struct {
	orderController controller.Order
}

func NewOrderHandler(orderController controller.Order) OrderHandler {
	return OrderHandler{
		orderController: orderController,
	}
}

func (c *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userID := handler.GetTokenClaim(r.Context()).ID
	orderID, err := uuid.GenerateUUID()
	if err != nil {
		handler.ResponseError(w, &handler.InternalServerError)
		return
	}

	rand.Seed(time.Now().UnixNano())
	max := 100000000

	randNum := rand.Intn(max)

	invoice := strconv.Itoa(randNum)
	invoice = "INV/" + invoice

	res, err := c.orderController.PlaceOrder(r.Context(), userID, orderID, invoice)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}

	handler.GenerateResponse(w, http.StatusOK, res)
}

func (c *OrderHandler) OrderLists(w http.ResponseWriter, r *http.Request) {
	userID := handler.GetTokenClaim(r.Context()).ID

	res, err := c.orderController.ListOfOrders(r.Context(), userID)
	if err != nil {
		handler.ResponseError(w, &handler.InternalServerError)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}

func (c *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	request := new(controller.UpdateOrderRequest)

	if err := handler.DecodeJSON(r, request); err != nil {
		handler.ResponseError(w, &handler.ParsePayloadError)
		return
	}

	request.Status = strings.ToLower(request.Status)
	err := c.orderController.UpdateOrderStatus(r.Context(), request)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}

func (c *OrderHandler) OrderDetail(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	err := validator.ValidateOrderID(orderID)
	var fieldError []handler.Fields
	if err != nil {
		fieldErr := handler.Fields{
			Name:    "order_id",
			Message: err.Error(),
		}
		fieldError = append(fieldError, fieldErr)
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldError)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := c.orderController.OrderDetail(r.Context(), orderID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}
