package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/gsc-lab/cs25-1-bannote-token-service/pkg/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTokenServiceServer
}

func (s *server) GenerateAccessToken(ctx context.Context, req *pb.GenerateAccessTokenRequest) (*pb.GenerateAccessTokenResponse, error) {
	log.Printf("GenerateAccessToken called: user_id=%s, claims=%v", req.UserId, req.Claims)

	// Mock response - echo back the input
	return &pb.GenerateAccessTokenResponse{
		AccessToken: fmt.Sprintf("mock_token_for_%s", req.UserId),
		ExpiresAt:   1234567890,
	}, nil
}

func (s *server) ValidateAccessToken(ctx context.Context, req *pb.ValidateAccessTokenRequest) (*pb.ValidateAccessTokenResponse, error) {
	log.Printf("ValidateAccessToken called: token=%s", req.AccessToken)

	// Mock response - always return valid
	return &pb.ValidateAccessTokenResponse{
		Valid:  true,
		UserId: "mock_user",
		Claims: map[string]string{"role": "admin"},
		Error:  "",
	}, nil
}

func main() {
	port := "9090"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTokenServiceServer(s, &server{})

	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
