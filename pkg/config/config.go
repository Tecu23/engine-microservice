package config

import (
	"os"
	"strings"
)

type Config struct {
	Port                int
	AuthTokens          []string
	TLSCertFile         string
	TLSKeyFile          string
	EnginePoolSize      int
	EnginePathStockfish string
}

func LoadConfig() *Config {
	port := 8089

	authTokensEnv := os.Getenv("AUTH_TOKENS")
	authTokens := strings.Split(authTokensEnv, ",")

	return &Config{
		Port:                port,
		AuthTokens:          authTokens,
		TLSCertFile:         os.Getenv("TLS_CERT_FILE"),
		TLSKeyFile:          os.Getenv("TLS_KEY_FILE"),
		EnginePoolSize:      4,
		EnginePathStockfish: "./bin/stockfish",
	}
}
