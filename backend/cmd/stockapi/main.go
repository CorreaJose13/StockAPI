package main

import (
	"context"
	"log"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/internal/analysis"
	"github.com/CorreaJose13/StockAPI/internal/api"
	"github.com/CorreaJose13/StockAPI/internal/db"
	"github.com/CorreaJose13/StockAPI/internal/repository"
	"github.com/CorreaJose13/StockAPI/models"
	"github.com/CorreaJose13/StockAPI/utils"
)

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

	analysis := analysis.NewAnalysis(formattedStocks)
	analysis.Analyze()

	if err := bulkInsert(ctx, analysis.Stocks); err != nil {
		log.Fatalf("failed to bulk insert stocks: %v", err)
	}
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

func bulkInsert(ctx context.Context, stocks []*models.FormattedStock) error {

	if err := repository.BulkInsertStocks(ctx, stocks); err != nil {
		return err
	}

	log.Printf("successfully inserted %d stocks into the database", len(stocks))

	return nil
}
