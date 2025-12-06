package middleware

import (
	dto "gin-backend-app/internal/dto/common"
	"gin-backend-app/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")

        if authHeader == "" {
			dto.SendError(c ,http.StatusUnauthorized, "Authorization header missing")
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			dto.SendError(c ,http.StatusUnauthorized, "Invalid authorization format. Use: Bearer <token>")
            return
        }

        tokenString := parts[1]

        // Validate token using utils
        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
			dto.SendError(c ,http.StatusUnauthorized, "Invalid or expired token")
            return
        }

        // Set user info in context
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("user_name", claims.Name)

        c.Next()
    }
}