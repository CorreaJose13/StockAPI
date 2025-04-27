package repository

import (
	"context"

	"github.com/CorreaJose13/StockAPI/models"
)

type StockRepository interface {
	BulkInsertStocks(ctx context.Context, stocks []*models.FormattedStock) error
	GetStocks(ctx context.Context) ([]*models.FormattedStock, error)
	Close() error
}

var stockRepoImpl StockRepository

func SetStockRepository(repo StockRepository) {
	stockRepoImpl = repo
}

func BulkInsertStocks(ctx context.Context, stocks []*models.FormattedStock) error {
	return stockRepoImpl.BulkInsertStocks(ctx, stocks)
}

func GetStocks(ctx context.Context) ([]*models.FormattedStock, error) {
	return stockRepoImpl.GetStocks(ctx)
}

func Close() error {
	return stockRepoImpl.Close()
}
