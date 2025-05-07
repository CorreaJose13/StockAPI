package repository

import (
	"context"

	"github.com/CorreaJose13/StockAPI/models"
)

type StockRepository interface {
	BulkInsertStocks(ctx context.Context, stocks []*models.FormattedStock) error
	BulkUpdateStocks(ctx context.Context, stocks []*models.FormattedStock) error
	GetStocks(ctx context.Context) ([]*models.FormattedStock, error)
	GetTableLength(ctx context.Context) (int, error)
	GetStocksFiltered(ctx context.Context, field, order, search string, page, limit int) ([]*models.FormattedStock, error)
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

func GetTableLength(ctx context.Context) (int, error) {
	return stockRepoImpl.GetTableLength(ctx)
}

func GetStocksFiltered(ctx context.Context, field, order, search string, page, limit int) ([]*models.FormattedStock, error) {
	return stockRepoImpl.GetStocksFiltered(ctx, field, order, search, page, limit)
}

func BulkUpdateStocks(ctx context.Context, stocks []*models.FormattedStock) error {
	return stockRepoImpl.BulkUpdateStocks(ctx, stocks)
}

func Close() error {
	return stockRepoImpl.Close()
}
