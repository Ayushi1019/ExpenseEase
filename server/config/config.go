package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config is a struct to hold the configuration variables
type Config struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	Port       string
	JwtSecret  string
}

// GetConfig is a function to get the configuration variables
func GetConfig() Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	config := Config{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASS"),
		DbName:     os.Getenv("DB_NAME"),
		JwtSecret:  os.Getenv("JWT_SECRET"),
	}

	return config
}
