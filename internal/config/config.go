package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL string
	BaseURL     string
	Secret      string
}

var C Config

func init() {

	_ = godotenv.Load()

	C = Config{
		DatabaseURL: getEnv("POSTGRES_URL", ""),
		BaseURL:     getEnv("BASE_URL", ""),
		Secret:      getEnv("SECRET_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
