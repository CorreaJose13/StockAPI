package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIURL      string
	BearerToken string
	DBURL       string
	APIKey      string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		APIURL:      os.Getenv("API_URL"),
		BearerToken: os.Getenv("BEARER_TOKEN"),
		DBURL:       os.Getenv("DB_URL"),
		APIKey:      os.Getenv("API_KEY"),
	}

	return config, nil
}
