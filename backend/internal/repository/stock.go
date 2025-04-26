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

var stockRepoImpl StockRepository

func SetStockRepository(repo StockRepository) {
	stockRepoImpl = repo
}

func InsertStock(ctx context.Context, stock *models.FormattedStock) error {
	return stockRepoImpl.InsertStock(ctx, stock)
}

func GetStocks(ctx context.Context) ([]*models.FormattedStock, error) {
	return stockRepoImpl.GetStocks(ctx)
}

func Close() error {
	return stockRepoImpl.Close()
}
