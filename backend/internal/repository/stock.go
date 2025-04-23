package repository

import (
	"context"

	"github.com/CorreaJose13/StockAPI/models"
)

type stock models.Response

type StockRepository interface {
	InsertStock(ctx context.Context, stock *stock) error
	GetStocks(ctx context.Context) ([]*stock, error)
	Close() error
}

var implementation StockRepository

func SetStockRepository(repo StockRepository) {
	implementation = repo
}

func InsertStock(ctx context.Context, stock *stock) error {
	return implementation.InsertStock(ctx, stock)
}

func GetStocks(ctx context.Context) ([]*stock, error) {
	return implementation.GetStocks(ctx)
}

func Close() error {
	return implementation.Close()
}
