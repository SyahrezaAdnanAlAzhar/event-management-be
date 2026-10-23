package config

import (
	"fmt"
	"os"
)

// Config holds all configuration values for the application
type Config struct {
	Port        string
	Env         string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	DatabaseURL string
}

// Load reads configuration from environment variables with sensible defaults
func Load() *Config {
	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "development")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "event_management")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			dbUser, dbPass, dbHost, dbPort, dbName, dbSSLMode)
	}

	return &Config{
		Port:        port,
		Env:         env,
		DBHost:      dbHost,
		DBPort:      dbPort,
		DBUser:      dbUser,
		DBPassword:  dbPass,
		DBName:      dbName,
		DBSSLMode:   dbSSLMode,
		DatabaseURL: dbURL,
	}
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
