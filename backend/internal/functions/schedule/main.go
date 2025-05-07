package main

import (
	"context"
	"log"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/internal/api"
	"github.com/CorreaJose13/StockAPI/internal/db"
	"github.com/CorreaJose13/StockAPI/internal/functions"
	"github.com/CorreaJose13/StockAPI/internal/repository"
	"github.com/CorreaJose13/StockAPI/models"
	"github.com/CorreaJose13/StockAPI/utils"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	repo    *db.CockRoachRepository
	cfg     *config.Config
	initErr error
)

func init() {
	repo, cfg, initErr = functions.FullSetup()
}

func handler(ctx context.Context) {
	if initErr != nil {
		log.Fatalf("failed to initialize: %v", initErr)
	}

	stocks := fetchStocks(cfg)

	formattedStocks := formatStocks(stocks)

	if err := repository.BulkUpdateStocks(context.Background(), formattedStocks); err != nil {
		log.Fatalf("failed to insert stocks: %v", err)
	}
}

func main() {
	lambda.Start(handler)
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

func formatStocks(stocks []models.Stock) []*models.FormattedStock {
	var formattedStocks []*models.FormattedStock
	for _, stock := range stocks {
		formattedStock, err := utils.Formatter(&stock)
		if err != nil {
			log.Fatalf("failed to format stock: %v", err)
		}
		formattedStocks = append(formattedStocks, formattedStock)
	}
	return formattedStocks
}
