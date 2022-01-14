package controller

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAuthController struct {
	mock.Mock
}

func (m *MockAuthController) Register(ctx context.Context, form *RegisterRequest) (string, error) {
	args := m.Called(ctx, form)
	return args.String(0), args.Error(1)
}

func (m *MockAuthController) ActivateEmail(ctx context.Context, email, id string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthController) Login(ctx context.Context, form *LoginRequest) (*LoginResponse, error) {
	args := m.Called(ctx, form)

	return &LoginResponse{}, args.Error(1)
}
