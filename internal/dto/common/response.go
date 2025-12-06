package common

import "github.com/gin-gonic/gin"

// Response represents standard API response
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Error message"`
}

// SendResponse sends success response
func SendResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends error response
func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
	})
}