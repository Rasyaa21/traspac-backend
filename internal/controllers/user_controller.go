package controllers

import (
	"gin-backend-app/internal/dto/common"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate user with email and password. Returns JWT token on successful authentication.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.LoginUserRequest true "Login credentials"
// @Success 200 {object} common.Response "Login successful"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 401 {object} common.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (uc *UserController) LoginUser(c *gin.Context) {
	var req request.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid Request Data")
		return
	}

	loginResponse, err := uc.UserService.LoginUser(req)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, loginResponse, "Login successful")
}

// RegisterUser godoc
// @Summary Register new user
// @Description Create a new user account with email verification. Returns JWT token and sends verification email.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "User registration data"
// @Success 201 {object} common.Response "User created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data or user already exists"
// @Router /auth/register [post]
func (uc *UserController) RegisterUser(c *gin.Context) {
	var req request.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid Request Data")
		return
	}

	registerResponse, err := uc.UserService.CreateUser(req)
	if err != nil {
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusCreated, registerResponse, "User created successfully, please check your email to use our feature")
}

// SendEmailPasswordReset godoc
// @Summary Request password reset
// @Description Send OTP code to user's email for password reset. OTP expires in 24 hours.
// @Tags Password Reset
// @Accept json
// @Produce json
// @Param request body request.RequestChangePasswordOtpRequest true "Email address for password reset"
// @Success 200 {object} common.Response "Password reset email sent successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 404 {object} common.ErrorResponse "Email not found"
// @Router /auth/request-change-password [post]
func (uc *UserController) SendEmailPasswordReset(c *gin.Context) {
	var req request.RequestChangePasswordOtpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid Request Data")
		return
	}

	err := uc.UserService.SendOtpToResetPassword(req.Email)
	if err != nil {
		common.SendError(c, http.StatusNotFound, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"email":   req.Email,
		"message": "Please check your email for verification code. Code expires in 24 Hour",
	}, "Password reset email sent successfully")
}

// GenerateAndSetVerificationToken godoc
// @Summary Verify OTP and get verification token
// @Description Verify OTP code received via email and get verification token for password reset
// @Tags Password Reset
// @Accept json
// @Produce json
// @Param request body request.VerifyOTPAndEmailRequest true "Email and OTP verification data"
// @Success 200 {object} common.Response "OTP verified successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data or OTP"
// @Failure 404 {object} common.ErrorResponse "Email not found or OTP expired"
// @Router /auth/verify-otp-password-change [post]
func (uc *UserController) GenerateAndSetVerificationToken(c *gin.Context) {
	var req request.VerifyOTPAndEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid Request Data")
		return
	}

	verificationToken, err := uc.UserService.GenerateAndSetVerificationToken(req.Email, req.TokenOtp)
	if err != nil {
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"verification_token": verificationToken,
		"message":           "Success validate token, use this verification token in params when changing password",
	}, "OTP verification successful")
}

// ValidateAndChangePassword godoc
// @Summary Change password with verification token
// @Description Change user password using verification token obtained from OTP verification
// @Tags Password Reset
// @Accept json
// @Produce json
// @Param token query string true "Verification token from OTP verification"
// @Param request body request.ChangePasswordRequest true "New password data"
// @Success 200 {object} common.Response "Password changed successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data or token"
// @Failure 404 {object} common.ErrorResponse "Invalid or expired verification token"
// @Router /auth/change-password [post]
func (uc *UserController) ValidateAndChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid Request Data")
		return
	}

	verificationToken := c.Query("token")
	if verificationToken == "" {
		common.SendError(c, http.StatusBadRequest, "Verification token is required")
		return
	}

	err := uc.UserService.ValidateAndChangePassword(verificationToken, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"message": "Password changed successfully",
	}, "Password change successful")
}

