package middleware

import (
	dto "gin-backend-app/internal/dto/common"
	"gin-backend-app/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, err := utils.GetClaimsFromHeader(c)
        if err != nil {
            dto.SendError(c, http.StatusUnauthorized, "Authentication required")
            c.Abort()
            return
        }

        if claims.UserID.String() == "" {
            dto.SendError(c, http.StatusUnauthorized, "Invalid user ID in token")
            c.Abort()
            return
        }

        c.Set("user", claims)
        c.Set("user_id", claims.UserID)        
        c.Set("user_email", claims.Email)      
        c.Set("user_name", claims.Name)  
        c.Set("verified_email", claims.VerifiedEmail)   

        c.Next()
    }
}