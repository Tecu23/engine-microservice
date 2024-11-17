package integration_test

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Tecu23/engine-microservice/pkg/api/generated"
	"github.com/Tecu23/engine-microservice/pkg/auth"
	"github.com/Tecu23/engine-microservice/pkg/server"
)

func TestAuthenticationInterceptor(t *testing.T) {
	// Setup test server
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	validToken := "test-token"
	auth.Initialize([]string{validToken})

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor()),
	)

	server.RegisterServer(grpcServer)

	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	// Setup test client
	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := generated.NewChessEngineClient(conn)

	tests := []struct {
		name       string
		token      string
		expectCode codes.Code
	}{
		{"Valid Token", validToken, codes.OK},
		{"Invalid Token", "invalid-token", codes.Unauthenticated},
		{"Missing Token", "", codes.Unauthenticated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.token != "" {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+tt.token)
			}
			_, err := client.CalculateBestMove(ctx, &generated.MoveRequest{})
			statusErr, _ := status.FromError(err)
			if statusErr.Code() != tt.expectCode {
				t.Errorf("Expected code: %v, got: %v", tt.expectCode, statusErr.Code())
			}
		})
	}
}
