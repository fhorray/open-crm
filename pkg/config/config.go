package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	DBHost string
	DBUser string
	DBPass string
}

var Cfg Config

func Load() {
	_ = godotenv.Load(".env")

	Cfg = Config{
		Port:   getEnv("PORT", "8787"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBUser: getEnv("DB_USER", "admin"),
		DBPass: getEnv("DB_PASS", "1234"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
