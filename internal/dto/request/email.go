package request

// EmailData represents email template data
type EmailData struct {
    Title   string `json:"title" example:"Verify Your Email Address"`
    Message string `json:"message" example:"Hi John, please verify your email to activate your account."`
    OTPCode string `json:"otp_code" example:"ABCD12"`
}

// OtpVerificationRequest represents OTP verification for email verification
type OtpVerificationRequest struct {
    TokenOtp string `json:"token_otp" validate:"required,len=6" example:"ABCD12" binding:"required,len=6"`
}