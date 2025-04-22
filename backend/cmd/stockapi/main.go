package main

import (
	"log"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/internal/api"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	consumer := api.NewApiConsumer(cfg)

	log.Println("Fetching stocks from API...")
	stocks, err := consumer.FetchStocks()
	if err != nil {
		log.Fatalf("Failed to fetch stocks: %v", err)
	}
	log.Printf("Successfully fetched %d stocks", len(stocks))
}
