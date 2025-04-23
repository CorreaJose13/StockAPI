package repository

import (
	"context"

	"github.com/CorreaJose13/StockAPI/models"
)

type StockRepository interface {
	InsertStock(ctx context.Context, stock *models.FormattedStock) error
	GetStocks(ctx context.Context) ([]*models.FormattedStock, error)
	Close() error
}

var implementation StockRepository

func SetStockRepository(repo StockRepository) {
	implementation = repo
}

func InsertStock(ctx context.Context, stock *models.FormattedStock) error {
	return implementation.InsertStock(ctx, stock)
}

func GetStocks(ctx context.Context) ([]*models.FormattedStock, error) {
	return implementation.GetStocks(ctx)
}

func Close() error {
	return implementation.Close()
}
