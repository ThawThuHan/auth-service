package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// EnvConfigLoader loads config from environment variables and .env file
type EnvConfigLoader struct {
	EnvPath string
}

// Load loads configuration from environment variables and .env file
func (l *EnvConfigLoader) Load() *Config {
	dotenvPath := l.EnvPath
	if dotenvPath == "" {
		dotenvPath = ".env"
		if _, err := os.Stat(dotenvPath); os.IsNotExist(err) {
			parentEnv := filepath.Join("..", "..", ".env")
			if _, err := os.Stat(parentEnv); err == nil {
				dotenvPath = parentEnv
			}
		}
	}
	if err := godotenv.Load(dotenvPath); err != nil {
		log.Printf("No .env file found at %s, relying on environment variables", dotenvPath)
	}
	cfg := &Config{}
	cfg.ENV = getEnv("ENV", "development")
	cfg.Host = getEnv("HOST", "localhost")
	cfg.Port = getEnv("PORT", "8080")
	cfg.JWTSecret = getEnv("JWT_SECRET", "")
	cfg.UserSvcAddr = getEnv("USER_SVC_ADDR", "localhost")
	cfg.UserSvcPort = getEnv("USER_SVC_PORT", "50051")
	cfg.PrivateKeyPath = getEnv("PRIVATE_KEY_PATH", "")
	cfg.PublicKeyPath = getEnv("PUBLIC_KEY_PATH", "")
	cfg.JWTKeyID = getEnv("JWT_KEY_ID", "")
	cfg.JWTIssuer = getEnv("JWT_ISSUER", "")
	cfg.JWTAudience = getEnv("JWT_AUDIENCE", "")
	return cfg
}

// getEnv returns the value of the environment variable or the default.
// If defaultValue is empty, the variable is considered required.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if defaultValue == "" {
		panic("Required environment variable not set: " + key)
	}
	return defaultValue
}
