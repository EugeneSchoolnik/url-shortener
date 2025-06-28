package handler

import "url-shortener/internal/service/auth"

type Dependencies struct {
	AuthService *auth.AuthService
}
