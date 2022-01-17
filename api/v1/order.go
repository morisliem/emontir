package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"e-montir/pkg/uuid"
	"net/http"
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

	res, err := c.orderController.PlaceOrder(r.Context(), userID, orderID)
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
