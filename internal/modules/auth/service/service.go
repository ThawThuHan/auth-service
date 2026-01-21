package service

import (
	"context"
	"errors"

	"github.com/TRIOSYS-Software/auth-service/internal/config"
	"github.com/TRIOSYS-Software/auth-service/internal/helper"
	"github.com/TRIOSYS-Software/auth-service/internal/modules/auth/dto"
	"github.com/TRIOSYS-Software/auth-service/internal/modules/auth/repository"
	"github.com/TRIOSYS-Software/auth-service/internal/util/jwt"
	"go.uber.org/zap"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
}

type authService struct {
	repository repository.AuthRepository
}

func NewAuthService(repository repository.AuthRepository) AuthService {
	return &authService{repository: repository}
}

func (s *authService) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	log := config.Cfg.Logger
	user, err := s.repository.VerifyUser(ctx, email)
	if err != nil {
		log.Error("Error verifying user", zap.Error(err))
		return nil, err
	}
	if user == nil {
		log.Warn("User not found", zap.String("email", email))
		return nil, errors.New("user not found")
	}
	if !helper.CheckPassword(user.Password, password) {
		log.Warn("Invalid password", zap.String("email", email))
		return nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.UserID, user.Username, user.Email)
	if err != nil {
		log.Error("Error generating token", zap.Error(err))
		return nil, err
	}

	return &dto.LoginResponse{Token: token}, nil
}
