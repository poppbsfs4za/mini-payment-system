package handlers

import (
	"errors"
	"net/http"
	"strings"

	"mini-payment-system/internal/domain/errs"

	"gorm.io/gorm"
)

type apiError struct {
	Status  int
	Code    string
	Message string
}

func mapError(err error) apiError {
	switch {
	case errors.Is(err, errs.ErrInvalidInput):
		return apiError{Status: http.StatusBadRequest, Code: "INVALID_INPUT", Message: err.Error()}
	case errors.Is(err, errs.ErrNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		return apiError{Status: http.StatusNotFound, Code: "NOT_FOUND", Message: err.Error()}
	case errors.Is(err, errs.ErrConflict):
		return apiError{Status: http.StatusConflict, Code: "CONFLICT", Message: err.Error()}
	case errors.Is(err, errs.ErrInsufficientBalance):
		return apiError{Status: http.StatusUnprocessableEntity, Code: "INSUFFICIENT_BALANCE", Message: err.Error()}
	case isDuplicateKeyError(err):
		return apiError{Status: http.StatusConflict, Code: "DUPLICATE_RESOURCE", Message: "resource already exists"}
	default:
		return apiError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "internal server error"}
	}
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key value") || strings.Contains(msg, "unique constraint")
}
