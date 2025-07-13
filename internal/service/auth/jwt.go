package auth

import (
	"time"
	"url-shortener/internal/service"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJWTService(secret string, duration time.Duration) *JWTService {
	return &JWTService{
		secretKey:     []byte(secret),
		tokenDuration: duration,
	}
}

type CustomClaims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func (s *JWTService) Generate(userID string) (string, error) {
	now := time.Now()

	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) Parse(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", service.ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", service.ErrInvalidToken
	}

	return claims.UserID, nil
}
