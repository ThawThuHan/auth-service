package main

import (
	"github.com/TRIOSYS-Software/auth-service/internal/config"
	grpcclient "github.com/TRIOSYS-Software/auth-service/internal/infrastructures/grpc-client"
	"github.com/TRIOSYS-Software/auth-service/internal/routes"
	"github.com/TRIOSYS-Software/auth-service/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize context with timeout
	config.Cfg.Logger.Info("Auth Service is starting...")

	// Initialize JWT Keys
	if err := jwt.InitializeJWT(
		config.Cfg.PrivateKeyPath,
		config.Cfg.PublicKeyPath,
		config.Cfg.JWTKeyID,
		config.Cfg.JWTIssuer,
		config.Cfg.JWTAudience,
	); err != nil {
		config.Cfg.Logger.Fatal("Failed to initialize JWT keys", zap.Error(err))
	}

	// Setup gRPC client
	grpcClient, err := grpcclient.NewGrpcClient(config.Cfg)
	if err != nil {
		config.Cfg.Logger.Fatal("Failed to create gRPC client", zap.Error(err))
	}
	defer func() {
		if err := grpcClient.Close(); err != nil {
			config.Cfg.Logger.Error("Failed to close gRPC client", zap.Error(err))
		}
	}()

	router := gin.New()

	authHandler := routes.ConfigureAuthHandler(grpcClient)
	routes.SetupRoutes(router, authHandler)

	if err := router.Run(config.Cfg.Host + ":" + config.Cfg.Port); err != nil {
		config.Cfg.Logger.Fatal("Failed to run HTTP server", zap.Error(err))
	}

	// Ensure logs are flushed to disk
	defer config.Cfg.Logger.Sync()
}
