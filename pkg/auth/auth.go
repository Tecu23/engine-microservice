package auth

import "errors"

var (
	ErrMissingToken = errors.New("missing authorization token")
	ErrInvalidToken = errors.New("invalid authorization token")
	validAuthTokens = map[string]bool{}
)

// Initialize valid tokens
func Initialize(validTokens []string) {
	for _, token := range validTokens {
		validAuthTokens[token] = true
	}
}

func ValidateToken(token string) error {
	if token == "" {
		return ErrMissingToken
	}

	if !validAuthTokens[token] {
		return ErrInvalidToken
	}

	return nil
}
