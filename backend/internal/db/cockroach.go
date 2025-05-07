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
	defaultOrder = "DESC"
	maxLimit     = 100
)

func ConnectCockRoachDB(cfg *config.Config) (*CockRoachRepository, error) {
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return &CockRoachRepository{db}, nil
}

func (repo *CockRoachRepository) GetStocksFiltered(ctx context.Context, field, order, search string, page, limit int) ([]*models.FormattedStock, error) {

	query, params := filterQueryParams(field, order, search, page, limit)
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

func (repo *CockRoachRepository) GetTableLength(ctx context.Context) (int, error) {
	var count int
	err := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stocks").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get table length: %w", err)
	}
	return count, nil
}

func (repo *CockRoachRepository) BulkInsertStocks(ctx context.Context, stocks []*models.FormattedStock) error {

	err := repo.createTable(ctx, "stocks")
	if err != nil {
		return fmt.Errorf("failed to create stocks table: %w", err)
	}

	err = repo.bulkInsertToTable(ctx, "stocks", stocks)
	if err != nil {
		return fmt.Errorf("failed to bulk insert stocks: %w", err)
	}

	return nil
}

func (repo *CockRoachRepository) BulkUpdateStocks(ctx context.Context, stocks []*models.FormattedStock) error {

	err := repo.createTable(ctx, "temp")
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	err = repo.bulkInsertToTable(ctx, "temp", stocks)
	if err != nil {
		return fmt.Errorf("failed to bulk insert stocks in temp table: %w", err)
	}

	count, err := repo.compareTables(ctx)
	if err != nil {
		return fmt.Errorf("failed to compare tables: %w", err)
	}

	if count > 0 {
		err = repo.mergeTables(ctx)
		if err != nil {
			return fmt.Errorf("failed to merge tables: %w", err)
		}

		err = repo.updateStockTable(ctx)
		if err != nil {
			return fmt.Errorf("failed to update stocks: %w", err)
		}
	}

	err = repo.DropTable(ctx, "temp")
	if err != nil {
		return fmt.Errorf("failed to drop temp table: %w", err)
	}

	return nil
}

func (repo *CockRoachRepository) compareTables(ctx context.Context) (int, error) {
	var count int
	err := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM (SELECT * FROM temp EXCEPT SELECT * FROM stocks);").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *CockRoachRepository) mergeTables(ctx context.Context) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		mergeQuery := `
        INSERT INTO stocks (ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time)
        SELECT t.ticker, t.target_from, t.target_to, t.company, t.action, t.brokerage, t.rating_from, t.rating_to, t.time
        FROM temp t
        LEFT JOIN stocks s ON t.ticker = s.ticker
        WHERE s.ticker IS NULL
    	`

		if _, err := tx.Exec(mergeQuery); err != nil {
			return err
		}

		return nil
	})
}

func (repo *CockRoachRepository) updateStockTable(ctx context.Context) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		updateQuery := `
		UPDATE stocks s
		SET
    		target_from = t.target_from,
    		target_to = t.target_to,
    		company = t.company,
    		time = t.time,
    		action = t.action,
    		brokerage = t.brokerage,
    		rating_from = t.rating_from,
    		rating_to = t.rating_to
		FROM temp t
		WHERE s.ticker = t.ticker
  			AND (s.target_from != t.target_from OR
       			s.target_to != t.target_to OR
       			s.time != t.time OR
       			s.company != t.company OR
       			s.action != t.action OR
       			s.brokerage != t.brokerage OR
       			s.rating_from != t.rating_from OR
       			s.rating_to != t.rating_to);
    	`

		if _, err := tx.Exec(updateQuery); err != nil {
			return err
		}

		return nil
	})
}

func (repo *CockRoachRepository) bulkInsertToTable(ctx context.Context, tableName string, stocks []*models.FormattedStock) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(pq.CopyIn(tableName, "ticker", "target_from", "target_to",
			"company", "action", "brokerage", "rating_from", "rating_to", "time"))
		if err != nil {
			return err
		}
		defer stmt.Close()

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

func (repo *CockRoachRepository) createTable(ctx context.Context, tableName string) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		ticker VARCHAR(10) PRIMARY KEY NOT NULL,
		target_from DECIMAL(10, 2) NOT NULL,
		target_to DECIMAL(10, 2) NOT NULL,
		company VARCHAR(100) NOT NULL,
		action VARCHAR(50) NOT NULL,
		brokerage VARCHAR(100) NOT NULL,
		rating_from VARCHAR(50) NOT NULL,
		rating_to VARCHAR(50) NOT NULL,
		time TIMESTAMP WITH TIME ZONE NOT NULL
		)`, pq.QuoteIdentifier(tableName))

		if _, err := tx.Exec(createTableQuery); err != nil {
			return err
		}

		return nil
	})
}

func (repo *CockRoachRepository) DropTable(ctx context.Context, tableName string) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		dropTableQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s", pq.QuoteIdentifier(tableName))
		if _, err := tx.Exec(dropTableQuery); err != nil {
			return err
		}
		return nil
	})
}

func filterQueryParams(field, order, search string, page, limit int) (string, []any) {
	page, limit = normalizePaginationParams(page, limit)
	offset := (page - 1) * limit
	query := generatePaginationQuery(field, order, search)
	params := make([]any, 0)
	searchParams := getSearchParams(search)
	params = append(params, searchParams...)
	params = append(params, limit, offset)

	return query, params
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

// validate max pages
func normalizePaginationParams(page, limit int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	if limit <= 0 {
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

func (repo *CockRoachRepository) execInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	return crdb.ExecuteTx(ctx, repo.db, nil, fn)
}

func (repo *CockRoachRepository) Close() error {
	return repo.db.Close()
}
