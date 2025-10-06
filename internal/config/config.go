package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
	JWT    JWTConfig
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret            string
	ExpirationMinutes int
}

func Load() (*Config, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "9090"),
		},
		JWT: JWTConfig{
			Secret:            jwtSecret,
			ExpirationMinutes: getEnvAsInt("JWT_EXPIRATION_MINUTES", 15),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
