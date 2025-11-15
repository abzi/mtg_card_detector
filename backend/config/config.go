package config

import (
	"os"
)

type Config struct {
	Port           string
	DatabasePath   string
	JWTSecret      string
	MigrationsPath string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabasePath:   getEnv("DATABASE_PATH", "./data/mtg_cards.db"),
		JWTSecret:      getEnv("JWT_SECRET", "change-this-in-production-to-a-secure-random-secret"),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "./migrations"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
