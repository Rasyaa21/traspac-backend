package services

import "gin-backend-app/internal/repositories"

type EmailVerficationService struct {
	userRepo *repositories.UserRepository
	userTokenEmail *repositories.UserTokenRepository
}

func NewEmailVerificationService(userRepo *repositories.UserRepository, userTokenEmailRepo *repositories.UserTokenRepository) *EmailVerficationService {
	return &EmailVerficationService{userRepo: userRepo, userTokenEmail: userTokenEmailRepo}
}

  