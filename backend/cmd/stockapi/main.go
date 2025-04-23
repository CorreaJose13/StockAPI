package main

import (
	"context"
	"log"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/internal/api"
	"github.com/CorreaJose13/StockAPI/internal/db"
	"github.com/CorreaJose13/StockAPI/internal/repository"
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

	defer repository.Close()

	consumer := api.NewApiConsumer(cfg)

	log.Println("fetching stocks from API...")
	stocks, err := consumer.FetchStocks()
	if err != nil {
		log.Fatalf("failed to fetch stocks: %v", err)
	}

	log.Printf("successfully fetched %d stocks", len(stocks))

	log.Println("Storing stocks in database...")
	for _, stock := range stocks {
		formattedStock, err := utils.Formatter(&stock)
		if err != nil {
			log.Printf("error formatting stock %s: %v", stock.Ticker, err)
		}

		if err := repository.InsertStock(ctx, formattedStock); err != nil {
			log.Printf("error storing stock %s: %v", stock.Ticker, err)
			continue
		}
	}

	log.Println("Successfully stored stocks in database")

	getStocks, err := repository.GetStocks(ctx)
	if err != nil {
		log.Fatalf("failed to get stocks: %v", err)
	}

	log.Printf("successfully retrieved %d stocks from database", len(getStocks))
}
