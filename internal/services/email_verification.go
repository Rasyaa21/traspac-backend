package services

import (
	"errors"
	"fmt"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/models"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/pkg/utils"
	"log"
	"time"

	"github.com/google/uuid"
)

type EmailVerficationService struct {
	UserRepo *repositories.UserRepository
	UserTokenEmail *repositories.UserTokenRepository
	Mailer *utils.SMTPMailer
	BaseUrl string
}

func NewEmailVerificationService(userRepo *repositories.UserRepository, userTokenEmailRepo *repositories.UserTokenRepository, mailer *utils.SMTPMailer, baseUrl string ) *EmailVerficationService {
	return &EmailVerficationService{UserRepo: userRepo, UserTokenEmail: userTokenEmailRepo, Mailer: mailer, BaseUrl: baseUrl }
}

func (s *EmailVerficationService) GetUserByID (userId uuid.UUID) (*models.User, error) {
	return s.UserRepo.FindByID(userId)
}

func (s *EmailVerficationService) VerifiyEmail (otp string, userId uuid.UUID) error {

	token, err := s.UserTokenEmail.FindTokenByOTP(otp, models.TokenTypeEmailVerification, userId)
	if err != nil {
		return errors.New("invalid otp")
	}

	if time.Now().After(token.ExpiresAt) {
		return errors.New("otp is already expired")
	}

	if token.UsedAt != nil {
		return errors.New("token is already used")
	}

	if err := s.UserRepo.MarkEmailVerified(token.UserID); err != nil {
		return err
	} 
	return s.UserTokenEmail.MarkAsUsedToken(token.ID)
}

func (s *EmailVerficationService) SendEmailVerification(user *models.User, emailVerificationType models.TokenType) error {
	otpToken, err := utils.GenerateUppercaseSixDigitOTP()
	if err != nil {
		return fmt.Errorf("failed to generate otp: %w", err)
	}

	if err := s.UserTokenEmail.Create(&models.UserToken{
		UserID:    user.ID,
		TokenOtp:  otpToken,
		TokenType: emailVerificationType,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}); err != nil {
		return fmt.Errorf("failed to store user token: %w", err)
	}

	emailData := &request.EmailData{
		OTPCode: otpToken,
	}

	var (
		subject string
		html    string
	)

	switch emailVerificationType {
	case models.TokenTypeEmailVerification:
		emailData.Title = "Verify Your Email Address"
		emailData.Message = fmt.Sprintf(
			"Hi %s, please verify your email to activate your account.",
			user.Name,
		)

		subject = "Email Verification"
		html, err = utils.BuildVerificationEmailHTML(emailData)

	case models.TokenTypePasswordReset:
		emailData.Title = "Reset Your Password"
		emailData.Message = fmt.Sprintf(
			"Hi %s, use the following code to reset your password.",
			user.Name,
		)

		subject = "Password Reset Request"
		html, err = utils.BuildResetPasswordEmailHTML(emailData)

	default:
		return fmt.Errorf("unsupported token type: %s", emailVerificationType)
	}

	if err != nil {
		return fmt.Errorf("failed to build email html: %w", err)
	}

	if err := s.Mailer.Send(user.Email, subject, html); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *EmailVerficationService) ResendEmailverification (user *models.User , emailVerificationType models.TokenType) error {
	if user.IsEmailVerified {
		return errors.New("email is already verified")
	}
	canResend := s.UserTokenEmail.CheckUserTokenCreatedAt(user.ID)
	if !canResend {
		return errors.New("you already request new verification link, try again in one minute")
	}
	return s.SendEmailVerification(user, models.TokenTypeEmailVerification)
}

func (s *EmailVerficationService) CleanupExpiredTokens() error {
	now := time.Now()
	deleted, err := s.UserTokenEmail.DeleteExpired(now)
	if err != nil {
		return err
	}
	log.Printf("[SERVICE] CleanupExpiredTokens: %d tokens deleted\n", deleted)
	return nil
}