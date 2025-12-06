package request

// CreateUserRequest represents user registration request
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required" example:"john_doe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"MyPass123!"`
}

// LoginUserRequest represents user login request
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required" example:"MyPass123!"`
}

// RequestPasswordChangeRequest represents request password change request
type RequestPasswordChangeRequest struct {
	Email     string    `json:"email" validate:"required,email"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	Password string `json:"password" validate:"required,min=8,regex=^(?=.*[A-Z])(?=.*[0-9])(?=.*[_!@#$%^&*])[A-Za-z0-9_!@#$%^&*]{8,}$"`
	NewPassword string `json:"new_password" validate:"required"`
}