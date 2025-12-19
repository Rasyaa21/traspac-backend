package middleware

import (
	dto "gin-backend-app/internal/dto/common"
	"gin-backend-app/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)
func RequireEmailVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		isVerified, err := utils.IsEmailVerifiedFromContext(c)
		if err != nil {
			dto.SendError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if !isVerified {
			dto.SendError(c, http.StatusForbidden, "Email not verified")
			c.Abort()
			return
		}

		c.Next()
	}
}