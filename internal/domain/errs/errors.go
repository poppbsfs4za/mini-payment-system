package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("conflict")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

func InvalidInput(message string) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, message)
}

func NotFound(message string) error {
	return fmt.Errorf("%w: %s", ErrNotFound, message)
}

func Conflict(message string) error {
	return fmt.Errorf("%w: %s", ErrConflict, message)
}
