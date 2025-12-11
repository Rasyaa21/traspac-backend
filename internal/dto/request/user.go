package request

// CreateUserRequest represents user registration request
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50" example:"john_doe" binding:"required"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com" binding:"required,email"`
	Password string `json:"password" validate:"required,min=8" example:"MySecurePass123!" binding:"required,min=8"`
}

// LoginUserRequest represents user login request
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com" binding:"required,email"`
	Password string `json:"password" validate:"required" example:"MySecurePass123!" binding:"required"`
}

// RequestChangePasswordOtpRequest represents request password reset request
type RequestChangePasswordOtpRequest struct {
	Email string `json:"email" validate:"required,email" example:"john@example.com" binding:"required,email"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"MyNewSecurePass123!" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required" example:"MyNewSecurePass123!" binding:"required"`
}

// VerifyOTPAndEmailRequest represents OTP verification request for password reset
type VerifyOTPAndEmailRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com" binding:"required,email"`
	TokenOtp string `json:"token_otp" validate:"required,len=6" example:"ABCD12" binding:"required,len=6"`
}