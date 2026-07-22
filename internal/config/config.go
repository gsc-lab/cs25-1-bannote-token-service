package config

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Server ServerConfig
	JWT    JWTConfig
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	PrivateKey        *rsa.PrivateKey
	PublicKey         *rsa.PublicKey
	ExpirationMinutes int
}

func Load() (*Config, error) {
	privateKeyBase64 := os.Getenv("PRIVATE_KEY_BASE64")
	publicKeyBase64 := os.Getenv("PUBLIC_KEY_BASE64")

	if privateKeyBase64 == "" {
		return nil, fmt.Errorf("PRIVATE_KEY_BASE64 environment variable is required")
	}

	if publicKeyBase64 == "" {
		return nil, fmt.Errorf("PUBLIC_KEY_BASE64 environment variable is required")
	}

	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)

	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)

	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	if err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "9091"),
		},
		JWT: JWTConfig{
			PrivateKey:        privateKey,
			PublicKey:         publicKey,
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
