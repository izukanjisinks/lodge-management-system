package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerPort string
	JWTSecret  string
	Email      EmailConfig
}

type EmailConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	FromEmail  string
	FromName   string
	UseTLS     bool
	UseSSL     bool
	AuthMethod string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "hr_system"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "hr-system-secret-key"),
		Email: EmailConfig{
			Host:       getEnv("EMAIL_HOST", "smtp.gmail.com"),
			Port:       getEnv("EMAIL_PORT", "587"),
			Username:   getEnv("EMAIL_USERNAME", ""),
			Password:   getEnv("EMAIL_PASSWORD", ""),
			FromEmail:  getEnv("EMAIL_FROM", ""),
			FromName:   getEnv("EMAIL_FROM_NAME", "HR System"),
			UseTLS:     getEnv("EMAIL_USE_TLS", "true") == "true",
			UseSSL:     getEnv("EMAIL_USE_SSL", "false") == "true",
			AuthMethod: getEnv("EMAIL_AUTH_METHOD", "PLAIN"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
