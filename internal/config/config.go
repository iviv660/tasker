package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL string
	BaseURL     string
}

var C Config

func init() {

	_ = godotenv.Load()

	C = Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		BaseURL:     getEnv("BASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
