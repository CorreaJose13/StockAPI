package db

import (
	"database/sql"
	"fmt"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/models"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

type CockRoachRepository struct {
	db *sql.DB
}

func NewPostgresRepository(cfg *config.Config) (*CockRoachRepository, error) {
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	createStocksTableQuery := `CREATE TABLE IF NOT EXISTS stocks (
        ticker VARCHAR(10) NOT NULL,
        target_from DECIMAL(10, 2) NOT NULL,
        target_to DECIMAL(10, 2) NOT NULL,
        company VARCHAR(255) NOT NULL,
        action VARCHAR(50) NOT NULL,
        brokerage VARCHAR(255) NOT NULL,
        rating_from VARCHAR(50) NOT NULL,
        rating_to VARCHAR(50) NOT NULL,
        time TIMESTAMP WITH TIME ZONE NOT NULL
    )`

	if _, err := db.Exec(createStocksTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create stocks table: %w", err)
	}

	return &CockRoachRepository{db}, nil
}

func (repo *CockRoachRepository) Close() error {
	return repo.db.Close()
}

func (repo *CockRoachRepository) InsertStock(ctx context.Context, stock *models.FormattedStock) error {
	insertStockQuery := `
        INSERT INTO stocks (
            ticker, target_from, target_to, company, 
            action, brokerage, rating_from, rating_to, time
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        )`

	_, err := repo.db.ExecContext(
		ctx,
		insertStockQuery,
		stock.Ticker,
		stock.TargetFrom,
		stock.TargetTo,
		stock.Company,
		stock.Action,
		stock.Brokerage,
		stock.RatingFrom,
		stock.RatingTo,
		stock.Time,
	)

	return err
}

func (repo *CockRoachRepository) GetStocks(ctx context.Context) ([]*models.FormattedStock, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM stocks")
	if err != nil {
		return nil, fmt.Errorf("failed to query stocks: %w", err)
	}
	defer rows.Close()

	var stocks []*models.FormattedStock
	for rows.Next() {
		var stock models.FormattedStock
		if err := rows.Scan(
			&stock.Id,
			&stock.Ticker,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Company,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
		); err != nil {
			return nil, fmt.Errorf("failed to scan stock: %w", err)
		}
		stocks = append(stocks, &stock)
	}

	return stocks, nil
}
