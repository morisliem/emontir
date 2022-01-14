package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"fmt"
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
	request.TotalPriceString = chi.URLParam(r, "total_price")

	fieldsErr, err := request.ValidatePaymentRequest()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := c.paymentController.Pay(r.Context(), request.OrderID, request.TotalPrice)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}

	handler.GenerateResponse(w, http.StatusOK, res)
}

// endpoint to check the payment status.
// it paid then order status will be updated and mechanic will be assigned to the order
func (c *PaymentHandler) PaymentNotification(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	request := new(controller.PaymentDetail)
	if err := handler.DecodeJSON(r, request); err != nil {
		handler.ResponseError(w, &handler.ParsePayloadError)
		return
	}

	fmt.Println("here")

	fmt.Println(request)

	if request.TransactionStatus == "PAID" {
		err := c.orderController.PaymentReceived(r.Context(), request.OrderID)
		if err != nil {
			handler.ResponseError(w, err)
		}
	}

	// token := chi.URLParam(r, "token")
	// err := validator.ValidatePaymentToken(token)
	// if err != nil {
	// 	var errField []handler.Fields
	// 	errField = append(errField, handler.Fields{
	// 		Name:    "payment_token",
	// 		Message: err.Error(),
	// 	})
	// 	res := handler.DefaultUnprocessableEntityError(handler.ValidationFailed, errField)
	// 	handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
	// 	return
	// }

	// res, err := c.paymentController.GetNotification(r.Context(), token)
	// if err != nil {
	// 	handler.ResponseError(w, err)
	// 	return
	// }

	// handler.GenerateResponse(w, http.StatusOK, res)
	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}
