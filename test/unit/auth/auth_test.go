package auth_test

import (
	"testing"

	"github.com/Tecu23/engine-microservice/pkg/auth"
)

func TestValidateToken(t *testing.T) {
	validTokens := []string{"token1", "token2"}
	auth.Initialize(validTokens)

	tests := []struct {
		name      string
		token     string
		expectErr bool
	}{
		{"Valid Token", "token1", false},
		{"Invalid Token", "invalid", true},
		{"Empty Token", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth.ValidateToken(tt.token)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}
