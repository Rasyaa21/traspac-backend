package services

import (
	"errors"
	"fmt"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/dto/response"
	"gin-backend-app/internal/models"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/pkg/utils"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepo *repositories.UserRepository
	UserTokenEmail *repositories.UserTokenRepository
	EmailService *EmailVerficationService
}

func NewUserService(userRepo *repositories.UserRepository, userTokenEmailRepo *repositories.UserTokenRepository, emailService *EmailVerficationService) *UserService {
    return &UserService{UserRepo: userRepo, UserTokenEmail: userTokenEmailRepo, EmailService: emailService}
}

func (s *UserService) CreateUser(req request.CreateUserRequest) (*response.LoginResponse, error) {
	email := req.Email
	name := req.Name
	existingUser, err := s.UserRepo.FindByEmailOrUsername(email, name)
	if err != nil {
		return nil ,err
	}

	if existingUser != nil {
		if existingUser.Email == email {
			return nil, errors.New("email already exists")
		}
		if existingUser.Name == name {
			return nil, errors.New("username already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hashed password")
	}
	user := models.User{
		Email: email,
		Name: name,
		Password: string(hashedPassword),
	}

	if err := s.UserRepo.Create(&user); err != nil {
        return nil, errors.New("failed to create user")
    }
	token, err := utils.CreateToken(&user)

	if err != nil {
		return nil, errors.New("token generation failed")
	}

	if s.EmailService != nil {
        if err := s.EmailService.SendEmailVerification(&user, models.TokenTypeEmailVerification); err != nil {
            log.Printf("Failed to send verification email: %v", err)
        }
    }

    return &response.LoginResponse{
		User: response.UserResponse{
			ID: user.ID,
			Name: user.Name,
			Email: user.Email,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
    }, nil
}

func (s *UserService) LoginUser(req request.LoginUserRequest) (*response.LoginResponse, error) {
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if user == nil {
        return nil, errors.New("user not found")
    }

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errors.New("wrong password")
	} 
	
	token, err := utils.CreateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return &response.LoginResponse{
        User: response.UserResponse{
            ID:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
        },
        Token: token,
    }, nil
}

func (s *UserService) SendOtpToResetPassword (email string) error {
	user, err := s.UserRepo.FindByEmail(email) 
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("email didnt exist")
		}
		return fmt.Errorf("failed to find user by email: %w", err)
	}

	if s.EmailService == nil {
		return errors.New("email service is not configured")
	}

	if err := s.EmailService.SendEmailVerification(user, models.TokenTypePasswordReset); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

func (s *UserService) GenerateAndSetVerificationToken(email, tokenOtp string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email) 
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "" ,errors.New("email didnt exist")
		}
		return "",fmt.Errorf("failed to find user by email: %w", err)
	}

	verifyToken, err := s.UserTokenEmail.GenerateAndSetVerificationTokenByOTP(tokenOtp, email, models.TokenTypePasswordReset, user.ID) 

	if err != nil {
		return "", errors.New(err.Error())
	}
	return verifyToken, nil
}


func (s *UserService) ValidateAndChangePassword (verificationToken, newPassword, confirmPassword string) error {
	user, err := s.UserTokenEmail.ValidateTokenAndGetUser(verificationToken)
	if err != nil {
		return err
	}

	if newPassword != confirmPassword {
		return errors.New("the new password and confirmation password must match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	err = s.UserRepo.ChangePassword(string(hashedPassword), user.ID)
	if err != nil {
		return err
	}
	return nil
}

