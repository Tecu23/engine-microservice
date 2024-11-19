package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port     int
	Version  string
	LogLevel string

	TLSCertFile string
	TLSKeyFile  string
}

type EngineConfig struct {
	PoolSize int
	Paths    map[string]string
}

type AuthConfig struct {
	AuthType   string
	AuthTokens []string
}

type Config struct {
	Server ServerConfig
	Engine EngineConfig
	Auth   AuthConfig
}

func InitConfig() (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("server.port", 8089)
	v.SetDefault("server.version", 8089)
	v.SetDefault("server.loglevel", "info")

	v.SetDefault("engine.poolsize", 5)

	v.SetDefault("auth.authtype", "apikey")

	// Allow overriding via environment variables
	v.SetEnvPrefix("CHESS_ENGINE")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read configuration from file (if any)
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/chess-engine/")

	// Read in the config file if it exists
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("fatal error config file: %w", err)
		}
	}

	// Unmarshal into the config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// EnginePathStockfish: "./bin/engines/stockfish",
	// EnginePathArgo:      "./bin/engines/argo",

	return &config, nil
}

func (cfg *Config) Validate() error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}
	if cfg.Engine.PoolSize <= 0 {
		return fmt.Errorf("engine pool size must be greater than 0")
	}
	// Add more validation rules as needed
	return nil
}
