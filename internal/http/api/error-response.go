package api

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrResponse(err string) ErrorResponse {
	return ErrorResponse{Error: err}
}
