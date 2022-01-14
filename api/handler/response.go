package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type (
	UnprocessableEntity struct {
		Message string   `json:"message"`
		Fields  []Fields `json:"fields"`
	}
	Fields struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
	DefaultError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	DefaultSuccess struct {
		Success bool `json:"success"`
	}
)

func GenerateResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonData, _ := json.Marshal(data)
	result := strings.Builder{}
	result.Write(jsonData)
	_, _ = w.Write([]byte(result.String()))
}

func RouteNotFound(w http.ResponseWriter, r *http.Request) {
	GenerateResponse(w, http.StatusNotFound, DefaultError{
		Message: "route not found",
	})
}

func (e *DefaultError) Error() string {
	return e.Message
}

func DefaultUnprocessableEntityError(msg string, fields []Fields) *UnprocessableEntity {
	return &UnprocessableEntity{
		Message: msg,
		Fields:  fields,
	}
}

func ResponseError(w http.ResponseWriter, err error) {
	var emontirErr *EmontirError
	if errors.As(err, &emontirErr) {
		msg := emontirErr.Message
		code := emontirErr.Code
		res := DefaultError{Message: msg, Code: code}
		if code == InternalServerError.Code {
			GenerateResponse(w, http.StatusInternalServerError, res)
			return
		}
		if code == ServiceNotExists.Code || code == CartAppointmentNotAvailable.Code || code == OrderNotExists.Code {
			GenerateResponse(w, http.StatusNotFound, res)
			return
		}
		GenerateResponse(w, http.StatusBadRequest, res)
		return
	}
	msg := InternalServerError.Message
	code := InternalServerError.Code
	res := DefaultError{Message: msg, Code: code}
	GenerateResponse(w, http.StatusInternalServerError, res)
}
