package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// auth
	ErrInvalidToken       = NewError(http.StatusUnauthorized, "invalid token")
	ErrInvalidCredentials = NewError(http.StatusBadRequest, "invalid credentials")
	// user
	ErrUserNotFound = NewError(http.StatusNotFound, "user not found")
	ErrEmailTaken   = NewError(http.StatusConflict, "email's already taken")
	// url
	ErrUrlNotFound      = NewError(http.StatusNotFound, "url not found")
	ErrAliasTaken       = NewError(http.StatusConflict, "this alias is already taken")
	ErrUrlStatsNotFound = NewError(http.StatusNotFound, "url statistics not found")
	// common
	ErrInternalError           = NewError(http.StatusInternalServerError, "internal server error")
	ErrValidation              = NewError(http.StatusBadRequest, "")
	ErrRelatedResourceNotFound = NewError(http.StatusUnprocessableEntity, "invalid entity")
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
