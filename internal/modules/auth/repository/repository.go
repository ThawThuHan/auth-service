package repository

import (
	"context"

	"github.com/TRIOSYS-Software/auth-service/internal/config"
	grpcclient "github.com/TRIOSYS-Software/auth-service/internal/infrastructures/grpc-client"
	"github.com/TRIOSYS-Software/auth-service/internal/modules/auth/dto"
	"go.uber.org/zap"
)

type AuthRepository interface {
	VerifyUser(ctx context.Context, email string) (*dto.GetUserByEmailResponse, error)
}

type authRepository struct {
	grpcClient grpcclient.GrpcClient
}

func NewAuthRepository(grpcClient grpcclient.GrpcClient) AuthRepository {
	return &authRepository{grpcClient: grpcClient}
}

func (r *authRepository) VerifyUser(ctx context.Context, email string) (*dto.GetUserByEmailResponse, error) {
	log := config.Cfg.Logger
	userResp, err := r.grpcClient.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("Failed to get user by email", zap.Error(err))
		return nil, err
	}
	if userResp.User == nil {
		log.Warn("User not found for email", zap.String("email", email))
		return nil, nil
	}

	return &dto.GetUserByEmailResponse{
		UserID:   userResp.User.Id,
		Email:    userResp.User.Email,
		Username: userResp.User.Username,
		Password: userResp.User.Password,
	}, nil
}
