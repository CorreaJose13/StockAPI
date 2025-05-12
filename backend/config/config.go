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
	APIKEY      string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := fullConfig()

	return config, nil
}

func LoadLambdaConfig() *Config {
	return fullConfig()
}

func LoadDbConfig() *Config {
	config := &Config{
		DBURL: os.Getenv("DB_URL"),
	}

	return config
}

func LoadAPIConfig() *Config {
	config := &Config{
		APIKEY: os.Getenv("API_KEY"),
	}

	return config
}

func fullConfig() *Config {
	config := &Config{
		APIURL:      os.Getenv("API_URL"),
		BearerToken: os.Getenv("BEARER_TOKEN"),
		DBURL:       os.Getenv("DB_URL"),
		APIKEY:      os.Getenv("API_KEY"),
	}

	return config
}
