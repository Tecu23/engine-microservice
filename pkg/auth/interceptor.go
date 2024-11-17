package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if info.FullMethod == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo" {
			return handler(ctx, req)
		}

		err := authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	var authHeader string

	if values := md["authorization"]; len(values) > 0 {
		authHeader = values[0]
	} else {
		return status.Error(codes.Unauthenticated, "authorization token not provided")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	token := parts[1]

	if err := ValidateToken(token); err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}

	return nil
}
