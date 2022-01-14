package v1

import (
	"e-montir/api/handler"
	"e-montir/controller"
	"e-montir/pkg/mailer"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authController   controller.Auth
	mailerController mailer.Mailer
}

func NewAuthHandler(authController controller.Auth, mailedCfg *mailer.Config) AuthHandler {
	return AuthHandler{
		authController:   authController,
		mailerController: *mailer.NewMailtrap(mailedCfg),
	}
}

func (c *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	request := new(controller.RegisterRequest)

	if err := handler.DecodeJSON(r, request); err != nil {
		handler.ResponseError(w, &handler.ParsePayloadError)
		return
	}

	fieldsErr, err := request.ValidateRegisterRequest()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	request.Email = strings.ToLower(request.Email)
	uid, err := c.authController.Register(r.Context(), request)

	if err != nil {
		handler.ResponseError(w, err)
		return
	}

	// case when user enter + in the email
	if strings.Contains(request.Email, "+") {
		request.Email = strings.ReplaceAll(request.Email, "+", "%2B")
	}

	go c.mailerController.SendActivationLink(r.Context(), request.Email, uid)
	handler.GenerateResponse(w, http.StatusOK, handler.DefaultSuccess{Success: true})
}

func (c *AuthHandler) ActivateEmail(w http.ResponseWriter, r *http.Request) {
	request := new(controller.ActivateEmailRequest)

	request.ID = r.URL.Query().Get("id")
	request.Email = strings.ToLower(r.URL.Query().Get("email"))

	fieldsErr, err := request.ValidateActivateEmailRequest()
	if err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	err = c.authController.ActivateEmail(r.Context(), request.Email, request.ID)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}

	http.Redirect(w, r, "emontir://login", http.StatusFound)
}

func (c *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	request := new(controller.LoginRequest)
	if err := handler.DecodeJSON(r, request); err != nil {
		handler.ResponseError(w, &handler.ParsePayloadError)
		return
	}

	if fieldsErr, err := request.ValidateLoginRequest(); err != nil {
		res := handler.DefaultUnprocessableEntityError(err.Error(), fieldsErr)
		handler.GenerateResponse(w, http.StatusUnprocessableEntity, res)
		return
	}

	request.Email = strings.ToLower(request.Email)
	token, err := c.authController.Login(r.Context(), request)
	if err != nil {
		handler.ResponseError(w, err)
		return
	}

	handler.GenerateResponse(w, http.StatusOK, token)
}
