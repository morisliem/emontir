package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"e-montir/pkg/uuid"
	"math/rand"
	"net/http"
	"strconv"
	"time"
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
