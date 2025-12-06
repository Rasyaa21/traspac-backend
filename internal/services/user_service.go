package services

import (
	"errors"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/dto/response"
	"gin-backend-app/internal/models"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req request.CreateUserRequest) (*response.LoginResponse, error) {
	email := req.Email
	name := req.Name
	existingUser, err := s.repo.FindByEmailOrUsername(email, name)
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

	if err := s.repo.Create(&user); err != nil {
        return nil, errors.New("failed to create user")
    }

	token, err := utils.CreateToken(&user)

	if err != nil {
		return nil, errors.New("token generation failed")
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
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if user == nil {
        return nil, errors.New("invalid credentials")
    }

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid credentials")
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


