package config

import "os"

type Config struct {
	AppEnv    string
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
}

func Load() Config {
	return Config{
		AppEnv:    getenv("APP_ENV", "development"),
		AppPort:   getenv("APP_PORT", "8080"),
		DBHost:    getenv("DB_HOST", "localhost"),
		DBPort:    getenv("DB_PORT", "5432"),
		DBUser:    getenv("DB_USER", "postgres"),
		DBPass:    getenv("DB_PASSWORD", "postgres"),
		DBName:    getenv("DB_NAME", "arise_assignment"),
		DBSSLMode: getenv("DB_SSLMODE", "disable"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
