package response

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse represents user data in response
type UserResponse struct {
    ID        uuid.UUID      `json:"id" example:"1"`
    Name      string    `json:"name" example:"john_doe"`
    Email     string    `json:"email" example:"john@example.com"`
    CreatedAt time.Time `json:"created_at"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
    User  UserResponse `json:"user"`
    Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
