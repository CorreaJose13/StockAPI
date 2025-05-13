package repository

import (
	"context"

	"github.com/CorreaJose13/StockAPI/models"
)

type StockRepository interface {
	BulkInsertStocks(ctx context.Context, stocks []*models.FormattedStock, tableName string) error
	BulkUpdateStocks(ctx context.Context, stocks []*models.FormattedStock, originalTable, tempTable string) error
	GetStocks(ctx context.Context, tableName string) ([]*models.FormattedStock, error)
	GetTableLength(ctx context.Context, tableName string) (int, error)
	GetStocksFiltered(ctx context.Context, field, order, search, tableName string, page, limit int) ([]*models.FormattedStock, error)
	Close() error
}

var stockRepoImpl StockRepository

func SetStockRepository(repo StockRepository) {
	stockRepoImpl = repo
}

func BulkInsertStocks(ctx context.Context, stocks []*models.FormattedStock, tableName string) error {
	return stockRepoImpl.BulkInsertStocks(ctx, stocks, tableName)
}

func BulkUpdateStocks(ctx context.Context, stocks []*models.FormattedStock, originalTable, tempTable string) error {
	return stockRepoImpl.BulkUpdateStocks(ctx, stocks, originalTable, tempTable)
}

func GetStocks(ctx context.Context, tableName string) ([]*models.FormattedStock, error) {
	return stockRepoImpl.GetStocks(ctx, tableName)
}

func GetTableLength(ctx context.Context, tableName string) (int, error) {
	return stockRepoImpl.GetTableLength(ctx, tableName)
}

func GetStocksFiltered(ctx context.Context, field, order, search, tableName string, page, limit int) ([]*models.FormattedStock, error) {
	return stockRepoImpl.GetStocksFiltered(ctx, field, order, search, tableName, page, limit)
}

func Close() error {
	return stockRepoImpl.Close()
}
