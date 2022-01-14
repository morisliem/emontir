package controller

import (
	"context"
	"database/sql"
	"e-montir/api/handler"
	"e-montir/model"
	"e-montir/pkg/jwt"
	"e-montir/pkg/password"
	"e-montir/pkg/uuid"
	"e-montir/pkg/validator"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type authCtx struct {
	userModel model.User
}

type Auth interface {
	Register(ctx context.Context, form *RegisterRequest) (string, error)
	ActivateEmail(ctx context.Context, email string, id string) error
	Login(ctx context.Context, form *LoginRequest) (*LoginResponse, error)
}

func NewAuth(userModel model.User) Auth {
	return &authCtx{
		userModel: userModel,
	}
}

type (
	RegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	ActivateEmailRequest struct {
		ID    string
		Email string
	}
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		Token     string `json:"token"`
		ExpiredAt string `json:"expired_at"`
	}
)

func (rr *RegisterRequest) ValidateRegisterRequest() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	err := validator.ValidateName(rr.Name)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "name",
			Message: err.Error(),
		})
	}
	err = validator.ValidateEmail(rr.Email)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "email",
			Message: err.Error(),
		})
	}
	err = validator.ValidatePassword(rr.Password)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "password",
			Message: err.Error(),
		})
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New("validation-failed")
}

func (v *ActivateEmailRequest) ValidateActivateEmailRequest() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	err := validator.ValidateID(v.ID)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "id",
			Message: err.Error(),
		})
	}
	err = validator.ValidateEmail(v.Email)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "email",
			Message: err.Error(),
		})
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

func (lr *LoginRequest) ValidateLoginRequest() ([]handler.Fields, error) {
	var count int
	var fields []handler.Fields
	err := validator.ValidateEmail(lr.Email)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "email",
			Message: err.Error(),
		})
	}

	err = validator.ValidatePassword(lr.Password)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "password",
			Message: "password wrong format",
		})
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

func (c *authCtx) Register(ctx context.Context, form *RegisterRequest) (string, error) {
	emailUsed, err := c.userModel.IsEmailUsed(ctx, form.Email)
	if err != nil {
		return "", &handler.InternalServerError
	}

	if emailUsed {
		return "", &handler.DuplicatedEmailError
	}

	uid, err := uuid.GenerateUUID()
	if err != nil {
		log.Error().Err(fmt.Errorf("error when generateUUID: %w", err)).Send()
		return "", &handler.InternalServerError
	}

	hashPassword, err := password.HashPassword(form.Password)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when hashPassword: %w", err)).Send()
		return "", &handler.InternalServerError
	}

	req := &model.RegisterUser{
		ID:       uid,
		Name:     form.Name,
		Email:    form.Email,
		Password: hashPassword,
	}

	err = c.userModel.RegisterUser(ctx, req)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when RegisterUser: %w", err)).Send()
		return "", &handler.InternalServerError
	}
	return uid, nil
}

func (c *authCtx) ActivateEmail(ctx context.Context, email, id string) error {
	res, err := c.userModel.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when GetUserByEmail: %w", err)).Send()
		if err == sql.ErrNoRows {
			return &handler.ActivationEmailFailedError
		}
		return &handler.InternalServerError
	}

	if res.ID != id {
		log.Error().Msg("incorrect id")
		return &handler.ActivationEmailFailedError
	}

	err = c.userModel.ActivateEmail(ctx, email)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when ActivateEmail: %w", err)).Send()
		return &handler.InternalServerError
	}
	return nil
}

func (c *authCtx) Login(ctx context.Context, form *LoginRequest) (*LoginResponse, error) {
	req := &model.LoginUser{
		Email:    form.Email,
		Password: form.Password,
	}

	res, err := c.userModel.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when GetUserByEmail: %w", err)).Send()
		if err == sql.ErrNoRows {
			return nil, &handler.LoginFailedError
		}
		return nil, &handler.InternalServerError
	}

	if !res.IsActive {
		log.Error().Msg("email not activated")
		return nil, &handler.EmailNotActivatedError
	}

	if !password.CompareHashPassword(req.Password, res.Password) {
		log.Error().Msg("incorrect password")
		return nil, &handler.LoginFailedError
	}

	keyDuration, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return nil, &handler.InternalServerError
	}

	token, expiredAt, err := jwt.GenerateToken(res.ID, os.Getenv("ACCESS_KEY"), keyDuration)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when generateAccessToken: %w", err)).Send()
		return nil, &handler.InternalServerError
	}

	return &LoginResponse{
		Token:     token,
		ExpiredAt: expiredAt,
	}, nil
}
