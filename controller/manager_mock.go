package controller

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockManagerController struct {
	mock.Mock
}

func (m *MockManagerController) Auth() Auth {
	MockAuthController := new(MockAuthController)

	MockAuthController.On("Register", context.Background(), &RegisterRequest{
		Name:     "hello world",
		Email:    "hello@gmail.com",
		Password: "abc123Dc1.",
	}).Return("", nil)

	MockAuthController.On("ActivateEmail", context.Background(), "hello@gmail.com").Return(nil)

	MockAuthController.On("Login", context.Background(), &LoginRequest{
		Email:    "hello@gmail.com",
		Password: "abc123Dc1.",
	}).Return(&LoginResponse{}, nil)

	return MockAuthController
}

func (m *MockManagerController) User() User {
	return nil
}

func (m *MockManagerController) Service() Service {
	return nil
}

func (m *MockManagerController) Timeslot() Timeslot {
	return nil
}

func (m *MockManagerController) Cart() Cart {
	return nil
}

func (m *MockManagerController) Order() Order {
	return nil
}

func (m *MockManagerController) Payment() Payment {
	return nil
}

func (m *MockManagerController) Review() Review {
	return nil
}
