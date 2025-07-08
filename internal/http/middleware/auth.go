package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JWTParser interface {
	Parse(tokenStr string) (string, error)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func AuthMiddleware(jwtParser JWTParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid authorization"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := jwtParser.Parse(tokenStr)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
			return
		}

		// Store in context
		c.Set("user_id", userID)
		c.Next()
	}
}
