package handler

import (
	"url-shortener/internal/service/auth"
	"url-shortener/internal/service/url"
)

type Dependencies struct {
	JwtService  *auth.JWTService
	AuthService *auth.AuthService
	UrlService  *url.UrlService
}
