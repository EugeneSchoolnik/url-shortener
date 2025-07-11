package api

import (
	"errors"
	"net/http"
	"url-shortener/internal/service"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrResponse(err string) ErrorResponse {
	return ErrorResponse{Error: err}
}

func ErrReponseFromServiceError(err error) (int, any) {
	statusCode := http.StatusInternalServerError
	var sErr *service.Error
	if errors.As(err, &sErr); sErr != nil {
		statusCode = sErr.StatusCode()
	}

	return statusCode, ErrResponse(err.Error())
}
