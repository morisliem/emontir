package v1

import (
	"bytes"
	"e-montir/controller"
	"e-montir/pkg/mailer"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var MockManagerController = new(controller.MockManagerController)
var route = GetHandler(MockManagerController, &mailer.Config{})

func TestRegister(t *testing.T) {
	tt := []struct {
		Name        string
		Input       *controller.RegisterRequest
		StatusCode  int
		ContentType string
	}{
		{
			Name: "Status ok",
			Input: &controller.RegisterRequest{
				Name:     "hello world",
				Email:    "hello@gmail.com",
				Password: "abc123Dc1.",
			},
			ContentType: "application/json",
			StatusCode:  http.StatusOK,
		},
		{
			Name: "Invalid param",
			Input: &controller.RegisterRequest{
				Name:     "hello world",
				Email:    "hello@gmail.com",
				Password: "abc123Dc1.",
			},
			StatusCode: http.StatusBadRequest,
		},
		{
			Name: "Unprocessable entity",
			Input: &controller.RegisterRequest{
				Name:     "hello world",
				Email:    "hello@gmailcom",
				Password: "abc123Dc1.",
			},
			ContentType: "application/json",
			StatusCode:  http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			rBody, err := json.Marshal(map[string]string{
				"name":     tc.Input.Name,
				"email":    tc.Input.Email,
				"password": tc.Input.Password,
			})
			if err != nil {
				t.Errorf("failed %v", err)
			}

			r := httptest.NewRequest("POST", "localhost:8080/api/v1/auth/register", bytes.NewBuffer(rBody))
			r.Header.Add("content-type", tc.ContentType)
			w := httptest.NewRecorder()

			route.Auth.Register(w, r)
			assert.Equal(t, w.Code, tc.StatusCode)

		})
	}
}

func TestActivateEmail(t *testing.T) {
	tt := []struct {
		Name       string
		Input      string
		StatusCode int
	}{
		{
			Name:       "Status found",
			Input:      "email=hello@gmail.com&id={user_id}",
			StatusCode: http.StatusFound,
		},
		{
			Name:       "URL param missing",
			Input:      "hello@gmail.com",
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name:       "Unprocessable entity",
			Input:      "email=hello@gmailcom",
			StatusCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "localhost:8080/api/v1/auth/verify?"+tc.Input, nil)
			w := httptest.NewRecorder()

			route.Auth.ActivateEmail(w, r)

			assert.Equal(t, w.Code, tc.StatusCode)
		})
	}
}

func TestLogin(t *testing.T) {
	tt := []struct {
		Name        string
		Input       *controller.LoginRequest
		StatusCode  int
		ContentType string
	}{
		{
			Name: "Status ok",
			Input: &controller.LoginRequest{
				Email:    "hello@gmail.com",
				Password: "abc123Dc1.",
			},
			ContentType: "application/json",
			StatusCode:  http.StatusOK,
		},
		{
			Name: "Invalid param",
			Input: &controller.LoginRequest{
				Email:    "hello@gmail.com",
				Password: "abc123Dc1.",
			},
			StatusCode: http.StatusBadRequest,
		},
		{
			Name: "Unprocessable entity",
			Input: &controller.LoginRequest{
				Email:    "hello@gmailcom",
				Password: "abc123Dc1.",
			},
			ContentType: "application/json",
			StatusCode:  http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			rBody, err := json.Marshal(map[string]string{
				"email":    tc.Input.Email,
				"password": tc.Input.Password,
			})
			if err != nil {
				t.Errorf("failed %v", err)
			}

			r := httptest.NewRequest("POST", "localhost:8080/api/v1/auth/login", bytes.NewBuffer(rBody))
			r.Header.Set("content-type", tc.ContentType)
			w := httptest.NewRecorder()

			route.Auth.Login(w, r)

			assert.Equal(t, w.Code, tc.StatusCode)
		})
	}
}
