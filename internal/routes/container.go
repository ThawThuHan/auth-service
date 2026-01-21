package routes

import (
	grpcclient "github.com/TRIOSYS-Software/auth-service/internal/infrastructures/grpc-client"
	"github.com/TRIOSYS-Software/auth-service/internal/modules/auth"
	"github.com/TRIOSYS-Software/auth-service/internal/modules/auth/repository"
	"github.com/TRIOSYS-Software/auth-service/internal/modules/auth/service"
)

func ConfigureAuthHandler(grpcclient grpcclient.GrpcClient) *auth.AuthHandler {
	repository := repository.NewAuthRepository(grpcclient)
	service := service.NewAuthService(repository)
	return auth.NewAuthHandler(service)
}
