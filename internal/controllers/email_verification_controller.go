package controllers

import (
	"gin-backend-app/internal/dto/common"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/models"
	"gin-backend-app/internal/services"
	"gin-backend-app/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailVerificationController struct {
	EmailVerificationService *services.EmailVerficationService
}

func NewEmailVerificationController(emailVerificationService *services.EmailVerficationService) *EmailVerificationController {
	return &EmailVerificationController{
		EmailVerificationService: emailVerificationService,
	}
}

// VerifyEmail godoc
// @Summary Verify email address with OTP
// @Description Verify user's email address using 6-digit OTP code received via email. User must be authenticated to use this endpoint.
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param request body request.OtpVerificationRequest true "OTP verification data"
// @Success 200 {object} common.Response "Email verified successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data or OTP"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "User not found"
// @Security BearerAuth
// @Router /auth/verify-email [post]
func (evc *EmailVerificationController) VerifyEmail(c *gin.Context) {
	var req request.OtpVerificationRequest

	log.Printf("DEBUG: Starting email verification")

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("DEBUG: GetUserIDFromContext error: %v", err)
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	log.Printf("DEBUG: UserID from context: %s", userID)

	user, err := evc.EmailVerificationService.GetUserByID(userID)
	if err != nil {
		log.Printf("DEBUG: GetUserByID error: %v", err)
		common.SendError(c, http.StatusNotFound, "User not found")
		return
	}

	log.Printf("DEBUG: User found: %s (%s)", user.Name, user.Email)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("DEBUG: JSON binding error: %v", err)
		common.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.TokenOtp == "" {
		log.Printf("DEBUG: Empty OTP token")
		common.SendError(c, http.StatusBadRequest, "OTP token is required")
		return
	}

	log.Printf("DEBUG: Attempting to verify OTP: %s", req.TokenOtp)

	if err := evc.EmailVerificationService.VerifiyEmail(req.TokenOtp, user.ID); err != nil {
		log.Printf("DEBUG: VerifiyEmail error: %v", err)
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("DEBUG: Email verification successful")
	common.SendResponse(c, http.StatusOK, gin.H{
		"email":   user.Email,
		"message": "Email verified successfully",
	}, "success")
}

// ResendEmailVerification godoc
// @Summary Resend email verification code
// @Description Resend 6-digit OTP verification code to user's email address. Rate limited to once per minute to prevent spam. User must be authenticated to use this endpoint.
// @Tags Email Verification
// @Accept json
// @Produce json
// @Success 200 {object} common.Response "Verification email sent successfully"
// @Failure 400 {object} common.ErrorResponse "Rate limit exceeded or email already verified"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "User not found"
// @Failure 429 {object} common.ErrorResponse "Too many requests - rate limited"
// @Security BearerAuth
// @Router /auth/resend-verification [post]
func (evc *EmailVerificationController) ResendEmailVerification(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	user, err := evc.EmailVerificationService.GetUserByID(userID)
	if err != nil {
		common.SendError(c, http.StatusNotFound, "User not found")
		return
	}

	if user.IsEmailVerified {
		common.SendError(c, http.StatusBadRequest, "Email is already verified")
		return
	}

	if err := evc.EmailVerificationService.SendEmailVerification(user, models.TokenTypeEmailVerification); err != nil {
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"email":   user.Email,
		"message": "Verification email sent successfully",
	}, "success")
}

// CheckVerificationStatus godoc
// @Summary Check email verification status
// @Description Get the current email verification status for the authenticated user
// @Tags Email Verification
// @Accept json
// @Produce json
// @Success 200 {object} common.Response "Verification status retrieved successfully"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "User not found"
// @Security BearerAuth
// @Router /auth/verification-status [get]
func (evc *EmailVerificationController) CheckVerificationStatus(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	user, err := evc.EmailVerificationService.GetUserByID(userID)
	if err != nil {
		common.SendError(c, http.StatusNotFound, "User not found")
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"user_id":             user.ID,
		"email":              user.Email,
		"is_email_verified":   user.IsEmailVerified,
		"email_verified_at":   user.EmailVerifiedAt,
	}, "Verification status retrieved successfully")
}