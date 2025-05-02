package db

import (
	"database/sql"
	"fmt"
	"strings"

	"context"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/models"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/lib/pq"
)

var (
	validFields = map[string]bool{
		"ticker":      true,
		"target_from": true,
		"target_to":   true,
		"company":     true,
		"action":      true,
		"brokerage":   true,
		"rating_from": true,
		"rating_to":   true,
		"time":        true,
	}

	validOrders = map[string]bool{
		"ASC":  true,
		"DESC": true,
	}
)

type CockRoachRepository struct {
	db *sql.DB
}

const (
	defaultPage  = 1
	defaultLimit = 10
	defaultField = "time"
	defaultOrder = "ASC"
	maxLimit     = 100
)

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
        time TIMESTAMP WITH TIME ZONE NOT NULL
    )`

	if _, err := db.Exec(createStocksTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create stocks table: %w", err)
	}

	return &CockRoachRepository{db}, nil
}

func (repo *CockRoachRepository) execInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	return crdb.ExecuteTx(ctx, repo.db, nil, fn)
}

func (repo *CockRoachRepository) Close() error {
	return repo.db.Close()
}

func (repo *CockRoachRepository) BulkInsertStocks(ctx context.Context, stocks []*models.FormattedStock) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(pq.CopyIn("stocks", "ticker", "target_from", "target_to",
			"company", "action", "brokerage", "rating_from", "rating_to", "time"))
		if err != nil {
			return err
		}
		defer stmt.Close()

		//formatTimeStamp
		for _, stock := range stocks {
			_, err = stmt.Exec(stock.Ticker, stock.TargetFrom, stock.TargetTo, stock.Company,
				stock.Action, stock.Brokerage, stock.RatingFrom, stock.RatingTo, stock.Time)
			if err != nil {
				return err
			}
		}

		_, err = stmt.Exec()
		if err != nil {
			return err
		}

		return nil
	})
}

func (repo *CockRoachRepository) GetStocks(ctx context.Context) ([]*models.FormattedStock, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM stocks")
	if err != nil {
		return nil, fmt.Errorf("failed to query stocks: %w", err)
	}
	defer rows.Close()

	stocks, err := scanRows(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan stocks: %w", err)
	}

	return stocks, nil
}

func (repo *CockRoachRepository) GetStocksFiltered(ctx context.Context, field, order, search string, page, limit int) ([]*models.FormattedStock, error) {
	page, limit = normalizePaginationParams(page, limit)
	offset := (page - 1) * limit

	query := generatePaginationQuery(field, order, search)
	params := make([]any, 0)
	searchParams := getSearchParams(search)
	params = append(params, searchParams...)
	params = append(params, limit, offset)

	rows, err := repo.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query stocks: %w", err)
	}
	defer rows.Close()

	stocks, err := scanRows(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan stocks: %w", err)
	}

	return stocks, nil
}

func (repo *CockRoachRepository) GetTableLength(ctx context.Context) (int, error) {
	var count int
	err := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stocks").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get table length: %w", err)
	}
	return count, nil
}

func generatePaginationQuery(field, order, search string) string {
	var pagination string
	baseQuery := "SELECT * FROM stocks"
	searchStm := buildSearchStatement(search)
	orderStm := buildOrderStatement(field, order)

	query := baseQuery
	if searchStm != "" {
		query += " " + searchStm
	}
	if orderStm != "" {
		query += " " + orderStm
	}

	if searchStm != "" {
		pagination = " LIMIT $2 OFFSET $3"
	} else {
		pagination = " LIMIT $1 OFFSET $2"
	}

	return query + pagination
}

func normalizePaginationParams(page, limit int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	if limit <= defaultLimit {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return page, limit
}

func isValidField(field string) bool {
	return validFields[strings.ToLower(field)]
}

func isValidOrder(order string) bool {
	return validOrders[strings.ToUpper(order)]
}

func buildSearchStatement(search string) string {
	if search == "" {
		return ""
	}
	return "WHERE ticker ILIKE $1 OR company ILIKE $1 OR brokerage ILIKE $1"
}

func getSearchParams(search string) []any {
	if search == "" {
		return []any{}
	}
	return []any{"%" + search + "%"}
}

func buildOrderStatement(field, order string) string {
	if !isValidField(field) {
		field = defaultField
	}

	if !isValidOrder(order) {
		order = defaultOrder
	}

	field = strings.ToLower(field)
	order = strings.ToUpper(order)

	return fmt.Sprintf("ORDER BY %s %s", field, order)
}

func scanRows(rows *sql.Rows) ([]*models.FormattedStock, error) {
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
			&stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}
	return stocks, nil
}
