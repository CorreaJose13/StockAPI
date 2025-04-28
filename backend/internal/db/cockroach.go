package db

import (
	"database/sql"
	"fmt"

	"context"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/models"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type CockRoachRepository struct {
	db *sql.DB
}

func NewPostgresRepository(cfg *config.Config) (*CockRoachRepository, error) {
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	createStocksTableQuery := `CREATE TABLE IF NOT EXISTS stocks (
        ticker VARCHAR(10) PRIMARY KEY NOT NULL,
        target_from DECIMAL(10, 2) NOT NULL,
        target_to DECIMAL(10, 2) NOT NULL,
        company VARCHAR(255) NOT NULL,
        action VARCHAR(20) NOT NULL,
        brokerage VARCHAR(255) NOT NULL,
        rating_from VARCHAR(50) NOT NULL,
        rating_to VARCHAR(50) NOT NULL,
		score DECIMAL(5, 2) NOT NULL,
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

func (repo *CockRoachRepository) BulkInsertStocks(ctx context.Context, stocks []*models.FormattedStock) error {
	txn, err := repo.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("stocks", "ticker", "target_from", "target_to",
		"company", "action", "brokerage", "rating_from", "rating_to", "time", "score"))
	if err != nil {
		return err
	}

	for _, stock := range stocks {
		_, err = stmt.Exec(stock.Ticker, stock.TargetFrom, stock.TargetTo, stock.Company,
			stock.Action, stock.Brokerage, stock.RatingFrom, stock.RatingTo, stock.Time, stock.Score)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
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
			&stock.Ticker,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Company,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Score,
			&stock.Time,
		); err != nil {
			return nil, fmt.Errorf("failed to scan stock: %w", err)
		}
		stocks = append(stocks, &stock)
	}

	return stocks, nil
}
