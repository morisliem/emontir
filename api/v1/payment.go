package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"net/http"

	"github.com/go-chi/chi"
)

type PaymentHandler struct {
	paymentController controller.Payment
	orderController   controller.Order
}

func NewPaymentHandler(paymentController controller.Payment, orderController controller.Order) PaymentHandler {
	return PaymentHandler{
		paymentController: paymentController,
		orderController:   orderController,
	}
}

// Endpoint to make payment as well as to store the checkout data to order table
func (c *PaymentHandler) MakePayment(w http.ResponseWriter, r *http.Request) {
	request := new(controller.PaymentRequest)
	request.OrderID = chi.URLParam(r, "order_id")

	fieldsErr, err := request.ValidatePaymentRequest()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := c.paymentController.Pay(r.Context(), request.OrderID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}

	handler.GenerateResponse(w, http.StatusOK, res)
}

// endpoint to check the payment status.
// it paid then order status will be updated and mechanic will be assigned to the order
func (c *PaymentHandler) PaymentNotification(w http.ResponseWriter, r *http.Request) {
	request := new(controller.PaymentDetail)
	if err := handler.DecodeJSON(r, request); err != nil {
		handler.ResponseError(w, &handler.ParsePayloadError)
		return
	}

	if request.TransactionStatus == "PAID" {
		err := c.orderController.PaymentReceived(r.Context(), request.OrderID)
		if err != nil {
			handler.ResponseError(w, err)
		}
	}

	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}
