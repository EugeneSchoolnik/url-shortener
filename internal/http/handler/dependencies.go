package handler

import (
	"url-shortener/internal/service/auth"
	clickstat "url-shortener/internal/service/click-stat"
	"url-shortener/internal/service/url"
)

type Dependencies struct {
	JwtService       *auth.JWTService
	AuthService      *auth.AuthService
	UrlService       *url.UrlService
	ClickStatService *clickstat.ClickStatService
}
