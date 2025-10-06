package main

import (
	"context"
	"fmt"
	"github.com/gsc-lab/cs25-1-bannote-token-service/pkg/jwt"
	"log"
	"net"

	pb "github.com/gsc-lab/cs25-1-bannote-token-service/pkg/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTokenServiceServer
	jwtManager *jwt.Manager
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
	port := "9090"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	jwtManager := jwt.NewManager("very-secret", 15)

	s := grpc.NewServer()
	pb.RegisterTokenServiceServer(s, &server{jwtManager: jwtManager})

	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
