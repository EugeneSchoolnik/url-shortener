package handler

import (
	"url-shortener/internal/service/auth"
	clickstat "url-shortener/internal/service/click-stat"
	"url-shortener/internal/service/url"
	"url-shortener/internal/service/user"
)

type Dependencies struct {
	JwtService       *auth.JWTService
	UserService      *user.UserService
	AuthService      *auth.AuthService
	UrlService       *url.UrlService
	ClickStatService *clickstat.ClickStatService
}
