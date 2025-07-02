package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// auth
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	// user
	ErrUserNotFound = errors.New("user not found")
	ErrEmailTaken   = errors.New("email's already taken")
	// common
	ErrInternalError = errors.New("internal server error")
	ErrValidation    = errors.New("")
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

	return fmt.Errorf("%w%s", ErrValidation, strings.Join(errMsgs, ", "))
}
