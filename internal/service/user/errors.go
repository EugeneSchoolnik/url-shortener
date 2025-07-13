package user

import "errors"

var (
	ErrEmailTaken         = errors.New("email's already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
