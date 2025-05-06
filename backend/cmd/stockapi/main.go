package main

import (
	"context"
	"log"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/internal/api"
	"github.com/CorreaJose13/StockAPI/internal/db"
	"github.com/CorreaJose13/StockAPI/internal/repository"
	"github.com/CorreaJose13/StockAPI/models"
	"github.com/CorreaJose13/StockAPI/utils"
)

// local function to fetch and store stocks retrieved from the API
func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	repo, err := db.NewPostgresRepository(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database repository: %v", err)
	}

	repository.SetStockRepository(repo)

	stocks := fetchStocks(cfg)

	var formattedStocks []*models.FormattedStock
	for _, stock := range stocks {
		formattedStock, err := utils.Formatter(&stock)
		if err != nil {
			log.Fatalf("failed to format stock: %v", err)
		}
		formattedStocks = append(formattedStocks, formattedStock)
	}

	if err := repository.BulkInsertStocks(ctx, formattedStocks); err != nil {
		log.Fatalf("failed to bulk insert stocks: %v", err)
	}

	log.Printf("successfully inserted %d stocks into the database", len(stocks))
}

func fetchStocks(cfg *config.Config) []models.Stock {
	consumer := api.NewAPIConsumer(cfg)

	log.Println("fetching stocks from API...")
	stocks, err := consumer.FetchStocks()
	if err != nil {
		log.Fatalf("failed to fetch stocks: %v", err)
	}

	log.Printf("successfully fetched %d stocks", len(stocks))

	return stocks
}
