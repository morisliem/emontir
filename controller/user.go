package controller

import (
	"context"
	"e-montir/api/handler"
	"e-montir/model"
	"e-montir/pkg/uuid"
	"e-montir/pkg/validator"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type userCtx struct {
	userModel model.User
}

type User interface {
	AddUserLocation(ctx context.Context, userID string, form *AddUserAddressRequest) error
	ListOfUserLocation(ctx context.Context, userID string) (*ListOfUserAddresses, error)
}

func NewUser(userModel model.User) User {
	return &userCtx{
		userModel: userModel,
	}
}

type (
	AddUserAddressRequest struct {
		Address       string `json:"address"`
		AddressDetail string `json:"address_detail"`
		Label         string `json:"label"`
		Recipient     string `json:"recipient"`
		PhoneNum      string `json:"phone_number"`
		Latitude      string `json:"latitude"`
		Longitude     string `json:"longitude"`
	}

	UserAddressResponse struct {
		ID            string `json:"id"`
		Address       string `json:"address"`
		AddressDetail string `json:"address_detail"`
		Label         string `json:"label"`
		Recipient     string `json:"recipient"`
		PhoneNum      string `json:"phone_number"`
		Latitude      string `json:"latitude"`
		Longitude     string `json:"longitude"`
	}

	UserCheckoutAddress struct {
		AddressID string `json:"id"`
		Address   string `json:"address"`
		Label     string `json:"label"`
		Recipient string `json:"recipient"`
		PhoneNum  string `json:"phone_number"`
	}

	ListOfUserAddresses struct {
		Address []UserAddressResponse `json:"addresses"`
	}
)

func (req *AddUserAddressRequest) ValidateAddUserLocation() ([]handler.Fields, error) {
	var fields []handler.Fields
	var count int

	err := validator.ValidateAddress(req.Address)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "address",
			Message: err.Error(),
		})
	}

	err = validator.ValidateLabel(req.Label)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "label",
			Message: err.Error(),
		})
	}

	err = validator.ValidateRecipientName(req.Recipient)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "recipient_name",
			Message: err.Error(),
		})
	}

	err = validator.ValidatePhoneNumber(req.PhoneNum)
	if err != nil {
		count++
		fields = append(fields, handler.Fields{
			Name:    "phone_number",
			Message: err.Error(),
		})
	}

	if count == 0 {
		return nil, nil
	}
	return fields, errors.New(handler.ValidationFailed)
}

func (c *userCtx) AddUserLocation(ctx context.Context, userID string, form *AddUserAddressRequest) error {
	locationID, err := uuid.GenerateUUID()
	if err != nil {
		return &handler.InternalServerError
	}

	err = c.userModel.AddUserLocation(ctx, userID, &model.UserLocation{
		ID:            locationID,
		Label:         form.Label,
		Address:       form.Address,
		AddressDetail: form.AddressDetail,
		PhoneNumber:   form.PhoneNum,
		RecipientName: form.Recipient,
		Latitude:      form.Latitude,
		Longitude:     form.Longitude,
	})

	if err != nil {
		log.Error().Err(fmt.Errorf("error when AddUserLocation: %w", err)).Send()
		return err
	}
	return nil
}

func (c *userCtx) ListOfUserLocation(ctx context.Context, userID string) (*ListOfUserAddresses, error) {
	listOfAddress := make([]UserAddressResponse, 0)
	res, err := c.userModel.GetListOfUserLocation(ctx, userID)
	if err != nil {
		log.Error().Err(fmt.Errorf("error when GetListOfUserLocation: %w", err)).Send()
		return nil, err
	}

	for i := 0; i < len(res); i++ {
		listOfAddress = append(listOfAddress, UserAddressResponse{
			ID:            res[i].ID,
			Label:         res[i].Label,
			Address:       res[i].Address,
			AddressDetail: res[i].AddressDetail,
			PhoneNum:      res[i].PhoneNumber,
			Recipient:     res[i].RecipientName,
			Latitude:      res[i].Latitude,
			Longitude:     res[i].Longitude,
		})
	}

	return &ListOfUserAddresses{
		Address: listOfAddress,
	}, nil
}
