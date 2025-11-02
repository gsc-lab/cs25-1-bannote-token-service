package main

import (
	"context"
	"fmt"
	"log"
	"net"

	healthpb "github.com/gsc-lab/cs25-1-bannote-token-service/gen/go/healthcheck"
	pb "github.com/gsc-lab/cs25-1-bannote-token-service/gen/go/token"
	"github.com/gsc-lab/cs25-1-bannote-token-service/internal/config"
	"github.com/gsc-lab/cs25-1-bannote-token-service/pkg/jwt"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTokenServiceServer
	jwtManager *jwt.Manager
}

type healthServer struct {
	healthpb.UnimplementedHealthServer
}

func (h *healthServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}

func (s *server) GenerateAccessToken(ctx context.Context, req *pb.GenerateAccessTokenRequest) (*pb.GenerateAccessTokenResponse, error) {
	log.Printf("GenerateAccessToken called: user_id=%s, roles=%s", req.UserId, req.Roles)

	if req.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	if req.Roles == "" {
		return nil, fmt.Errorf("roles is required")
	}

	token, expiresAt, err := s.jwtManager.GenerateToken(req.UserId, req.Roles)

	if err != nil {
		return nil, err
	}

	return &pb.GenerateAccessTokenResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *server) ValidateAccessToken(ctx context.Context, req *pb.ValidateAccessTokenRequest) (*pb.ValidateAccessTokenResponse, error) {
	log.Printf("ValidateAccessToken called: token=%s", req.AccessToken)

	claims, err := s.jwtManager.ValidateToken(req.AccessToken)

	if err != nil {
		return &pb.ValidateAccessTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &pb.ValidateAccessTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Roles:  claims.Roles,
		Error:  "",
	}, nil
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	jwtManager := jwt.NewManager(cfg.JWT.PrivateKey, cfg.JWT.PublicKey, cfg.JWT.ExpirationMinutes)

	s := grpc.NewServer()
	pb.RegisterTokenServiceServer(s, &server{jwtManager: jwtManager})
	healthpb.RegisterHealthServer(s, &healthServer{})

	log.Printf("gRPC server listening on port %s", cfg.Server.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
