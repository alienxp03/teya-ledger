package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ServiceError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorCode string

const (
	NotFound                 ErrorCode = "NOT_FOUND"
	BadRequest               ErrorCode = "BAD_REQUEST"
	ErrorCodeInvalidAmount   ErrorCode = "INVALID_AMOUNT"
	ErrorCodeInvalidCurrency ErrorCode = "INVALID_CURRENCY"
	ErrorInvalidParams       ErrorCode = "INVALID_PARAMS"
)

func (e ServiceError) Error() string {
	return e.Message
}

func (s *ServiceError) RespondError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(s.Status)
	if err := json.NewEncoder(w).Encode(s.Message); err != nil {
		fmt.Printf("Could not encode JSON body: %s\n", err.Error())
	}
}

func BadRequestError(code ErrorCode, message string) ServiceError {
	return ServiceError{
		Status:  http.StatusBadRequest,
		Code:    string(code),
		Message: message,
	}
}

var (
	NewBadRequest = func(code ErrorCode, message string) *ServiceError {
		return &ServiceError{
			Status:  http.StatusBadRequest,
			Code:    string(code),
			Message: message,
		}
	}

	NewNotFound = func(message string) *ServiceError {
		return &ServiceError{
			Status:  http.StatusNotFound,
			Code:    string(NotFound),
			Message: message,
		}
	}
)
