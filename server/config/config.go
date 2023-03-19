package config

import (
	"os"
)

// Config is a struct that holds the configuration for the application
type Config struct {
	Port   string
	DbHost string
	DbPort string
	DbUser string
	DbPass string
	DbName string
}

// GetConfig returns the configuration for the application
func GetConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	return Config{
		Port:   port,
		DbHost: os.Getenv("DB_HOST"),
		DbPort: os.Getenv("DB_PORT"),
		DbUser: os.Getenv("DB_USER"),
		DbPass: os.Getenv("DB_PASS"),
		DbName: os.Getenv("DB_NAME"),
	}
}
