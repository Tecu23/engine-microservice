package auth_test

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Tecu23/engine-microservice/pkg/auth"
	"github.com/Tecu23/engine-microservice/pkg/config"
)

func TestAuthenticate_ValidAPIKey(t *testing.T) {
	cfg := &config.AuthConfig{
		AuthType:   "apikey",
		AuthTokens: []string{"valid-api-key"},
	}
	auth.Initialize(cfg)

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("authorization", "ApiKey laXnteaAiKlPnCqJ"),
	)

	err := auth.Authenticate(ctx)
	if err != nil {
		t.Errorf("Expected authentication to succeed, got error: %v", err)
	}
}

func TestAuthenticate_InvalidAPIKey(t *testing.T) {
	cfg := &config.AuthConfig{
		AuthType:   "apikey",
		AuthTokens: []string{"valid-api-key"},
	}
	auth.Initialize(cfg)

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("authorization", "ApiKey invalid-api-key"),
	)

	err := auth.Authenticate(ctx)
	if err == nil {
		t.Error("Expected authentication to fail for invalid API key")
	} else {
		st, _ := status.FromError(err)
		if st.Code() != codes.Unauthenticated {
			t.Errorf("Expected Unauthenticated error, got %v", st.Code())
		}
	}
}

func TestAuthenticate_MissingAPIKey(t *testing.T) {
	cfg := &config.AuthConfig{
		AuthType:   "apikey",
		AuthTokens: []string{"valid-api-key"},
	}
	auth.Initialize(cfg)

	ctx := context.Background() // No metadata

	err := auth.Authenticate(ctx)
	if err == nil {
		t.Error("Expected authentication to fail when API key is missing")
	} else {
		st, _ := status.FromError(err)
		if st.Code() != codes.Unauthenticated {
			t.Errorf("Expected Unauthenticated error, got %v", st.Code())
		}
	}
}
