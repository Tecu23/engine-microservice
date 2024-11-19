package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Tecu23/engine-microservice/pkg/config"
)

var (
	ErrMissingToken = errors.New("missing authorization token")
	ErrInvalidToken = errors.New("invalid authorization token")
)

var validAPIKeys map[string]bool

// Initialize valid tokens
func Initialize(cfg *config.AuthConfig) error {
	if cfg.AuthType != "apikey" {
		return fmt.Errorf("unsupported auth type: %s", cfg.AuthType)
	}

	validAPIKeys = make(map[string]bool)
	for _, key := range cfg.AuthTokens {
		validAPIKeys[key] = true
	}

	return nil
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo" {
			return handler(ctx, req)
		}

		// Authenticate the request
		err := Authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func Authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	var apiKey string
	if values := md["authorization"]; len(values) > 0 {
		// Expecting "ApiKey <key>"
		parts := strings.SplitN(values[0], " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "apikey" {
			apiKey = parts[1]
		}
	}

	if apiKey == "" {
		return status.Error(codes.Unauthenticated, "API key not provided")
	}

	if !validAPIKeys[apiKey] {
		return status.Error(codes.Unauthenticated, "invalid API key")
	}

	return nil
}
