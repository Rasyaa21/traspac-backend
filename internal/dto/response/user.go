package response

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse represents user data in API responses
type UserResponse struct {
	ID               uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name             string     `json:"name" example:"john_doe"`
	Email            string     `json:"email" example:"john@example.com"`
	IsEmailVerified  bool       `json:"is_email_verified" example:"false"`
	EmailVerifiedAt  *time.Time `json:"email_verified_at,omitempty" example:"2023-01-01T00:00:00Z"`
	CreatedAt        time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt        time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// LoginResponse represents login/register response
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
