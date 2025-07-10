package api

import (
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
	if sErr, ok := err.(*service.Error); ok {
		statusCode = sErr.StatusCode()
	}

	return statusCode, ErrResponse(err.Error())
}
