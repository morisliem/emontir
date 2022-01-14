package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"net/http"
)

type UserHandler struct {
	userController controller.User
}

func NewUserHandler(userController controller.User) UserHandler {
	return UserHandler{
		userController: userController,
	}
}

func (c *UserHandler) AddUserLocation(w http.ResponseWriter, r *http.Request) {
	request := new(controller.AddUserAddressRequest)

	if err := handler.DecodeJSON(r, request); err != nil {
		handler.ResponseError(w, &handler.ParsePayloadError)
		return
	}

	fieldsErr, err := request.ValidateAddUserLocation()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	userID := handler.GetTokenClaim(r.Context()).ID
	err = c.userController.AddUserLocation(r.Context(), userID, request)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}

func (c *UserHandler) ListOfUserLocation(w http.ResponseWriter, r *http.Request) {
	userID := handler.GetTokenClaim(r.Context()).ID
	res, err := c.userController.ListOfUserLocation(r.Context(), userID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}
	handler.GenerateResponse(w, http.StatusOK, res)
}
