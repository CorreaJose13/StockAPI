package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiUrl      string
	BearerToken string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	config := &Config{
		ApiUrl:      os.Getenv("API_URL"),
		BearerToken: os.Getenv("BEARER_TOKEN"),
	}

	return config, nil
}
