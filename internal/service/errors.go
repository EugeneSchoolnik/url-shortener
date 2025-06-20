package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInternalError = errors.New("internal server error")
	ErrValidation    = errors.New("validation error")
	ErrNotFound      = errors.New("not found")
)

func PrettyValidationError(validationErrs validator.ValidationErrors) error {
	var errMsgs []string

	for _, err := range validationErrs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid email", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return fmt.Errorf("%w: %s", ErrValidation, strings.Join(errMsgs, ", "))
}
