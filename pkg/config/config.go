package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DB_PORT string
	DB_HOST string
	DB_USER string
	DB_PASS string
	DB_NAME string
}

type KeysConfig struct {
	JWT_SECRET   string
	STRIPE_KEY   string
	RESEND_TOKEN string
	// Adicione outras chaves sensíveis aqui
}

type AuthConfig struct {
	JWT_EXPIRES_IN               time.Duration
	JWT_REFRESH_TOKEN_EXPIRES_IN time.Duration
	JWT_ISSUER                   string
	JWT_AUDIENCE                 string
	DEFAULT_ROLE                 string
}

type Config struct {
	DB   DBConfig
	KEYS KeysConfig
	AUTH AuthConfig
}

var Cfg Config

func Load() {
	_ = godotenv.Load(".env")

	Cfg = Config{
		DB: DBConfig{
			DB_PORT: getEnv("PORT", "5432"),
			DB_HOST: getEnv("DB_HOST", "localhost"),
			DB_USER: getEnv("DB_USER", "admin"),
			DB_PASS: getEnv("DB_PASS", "1234"),
			DB_NAME: getEnv("DB_NAME", "postgres"),
		},
		KEYS: KeysConfig{
			JWT_SECRET:   getEnv("JWT_SECRET", "default-secret"),
			STRIPE_KEY:   getEnv("STRIPE_KEY", ""),
			RESEND_TOKEN: getEnv("RESEND_TOKEN", ""),
		},
		AUTH: AuthConfig{
			JWT_EXPIRES_IN:               mustParseDuration(getEnv("JWT_EXPIRES_IN", "24h")),
			JWT_REFRESH_TOKEN_EXPIRES_IN: mustParseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRES_IN", "168h")),
			JWT_ISSUER:                   getEnv("JWT_ISSUER", "api"),
			JWT_AUDIENCE:                 getEnv("JWT_AUDIENCE", "app"),
			DEFAULT_ROLE:                 "user",
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func mustParseDuration(value string) time.Duration {
	dur, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Invalid duration for key: %s - %v", value, err)
	}
	return dur
}
