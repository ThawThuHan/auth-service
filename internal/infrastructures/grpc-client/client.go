package grpcclient

import (
	"context"

	"github.com/TRIOSYS-Software/auth-service/internal/config"
	userv1 "github.com/TRIOSYS-Software/auth-service/proto/gen/proto/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient interface {
	ClientConn() *grpc.ClientConn
	Close() error
	GetUserByEmail(ctx context.Context, email string) (*userv1.GetUserByEmailResponse, error)
}

type grpcClient struct {
	conn *grpc.ClientConn
	userv1.UserServiceClient
}

func NewGrpcClient(cfg *config.Config) (GrpcClient, error) {
	conn, err := grpc.NewClient(cfg.UserSvcAddr+":"+cfg.UserSvcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cfg.Logger.Error("Failed to connect to User Service gRPC server", zap.Error(err))
		return nil, err
	}
	return &grpcClient{
		conn: conn,
	}, nil
}

func (g *grpcClient) ClientConn() *grpc.ClientConn {
	return g.conn
}

func (g *grpcClient) Close() error {
	return g.conn.Close()
}

func (g *grpcClient) GetUserByEmail(ctx context.Context, email string) (*userv1.GetUserByEmailResponse, error) {
	client := userv1.NewUserServiceClient(g.ClientConn())
	req := &userv1.GetUserByEmailRequest{Email: email}
	return client.GetUserByEmail(ctx, req)
}
