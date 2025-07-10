package service

type Error struct {
	statusCode int
	message    string
}

func NewError(statusCode int, message string) *Error {
	return &Error{statusCode, message}
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) StatusCode() int {
	return e.statusCode
}
