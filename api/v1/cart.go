package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"e-montir/pkg/validator"
	"net/http"

	"github.com/go-chi/chi"
)

type CartHandler struct {
	cartController controller.Cart
}

func NewCartHandler(cartController controller.Cart) CartHandler {
	return CartHandler{
		cartController: cartController,
	}
}

func (c *CartHandler) SetCartAppointment(w http.ResponseWriter, r *http.Request) {
	request := new(controller.CartAppointmentRequest)

	if err := handler.DecodeJSON(r, request); err != nil {
		handler.ResponseError(w, &handler.ParsePayloadError)
		return
	}

	request.UserID = handler.GetTokenClaim(r.Context()).ID
	fieldsErr, err := request.ValidateCartAppointment()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	err = c.cartController.SetCartAppointment(r.Context(), request.UserID, request.Date, request.Time)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}

func (c *CartHandler) RemoveCartAppointment(w http.ResponseWriter, r *http.Request) {
	userID := handler.GetTokenClaim(r.Context()).ID
	err := validator.ValidateID(userID)
	if err != nil {
		var errField []handler.Fields
		errField = append(errField, handler.Fields{
			Name:    "user_id",
			Message: err.Error(),
		})
		res := handler.DefaultUnprocessableEntityError(handler.ValidationFailed, errField)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	// Remove cart appointment will remove the entiry cart
	// because user and cart has one to one relationship
	// thus when user change date or timeslot appointment,
	// entiry data in the cart will be removed
	err = c.cartController.RemoveCartAppointment(r.Context(), userID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}

func (c *CartHandler) AddServiceToCart(w http.ResponseWriter, r *http.Request) {
	request := new(controller.AddOrRemoveServiceToCartRequest)
	request.CartIDString = handler.GetTokenClaim(r.Context()).ID

	if err := handler.DecodeJSON(r, request); err != nil {
		handler.ResponseError(w, &handler.ParsePayloadError)
		return
	}

	fieldsErr, err := request.ValidateAddOrRemoveServiceToCart()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := c.cartController.AddServiceToCart(r.Context(), request.ServiceID, request.CartIDString)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}

func (c *CartHandler) RemoveServiceFromCart(w http.ResponseWriter, r *http.Request) {
	request := new(controller.AddOrRemoveServiceToCartRequest)
	request.ServiceIDString = chi.URLParam(r, "service_id")
	userID := handler.GetTokenClaim(r.Context()).ID

	fieldsErr, err := request.ValidateAddOrRemoveServiceToCart()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := c.cartController.RemoveServiceFromCart(r.Context(), request.ServiceID, userID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}

func (c *CartHandler) GetCheckoutDetail(w http.ResponseWriter, r *http.Request) {
	userID := handler.GetTokenClaim(r.Context()).ID
	res, err := c.cartController.CartDetail(r.Context(), userID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}
