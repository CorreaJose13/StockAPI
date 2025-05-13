package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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

	ErrInvalidField = errors.New("invalid field")
	ErrInvalidOrder = errors.New("invalid order")
)

type CockRoachRepository struct {
	db *sql.DB
}

type queryBuilder struct {
	query  string
	params []any
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
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return &CockRoachRepository{db}, nil
}

func (repo *CockRoachRepository) GetStocksFiltered(ctx context.Context, field, order, search, tableName string, page, limit int) ([]*models.FormattedStock, error) {

	result, err := filterQueryParams(field, order, search, tableName, page, limit)
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryContext(ctx, result.query, result.params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query filtered stocks: %w", err)
	}

	defer rows.Close()

	stocks, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

func (repo *CockRoachRepository) GetStocks(ctx context.Context, tableName string) ([]*models.FormattedStock, error) {
	query := fmt.Sprintf(`SELECT * FROM %s`, pq.QuoteIdentifier(tableName))
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stocks: %w", err)
	}

	defer rows.Close()

	stocks, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

func (repo *CockRoachRepository) GetTableLength(ctx context.Context, tableName string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", pq.QuoteIdentifier(tableName))
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get table %s length: %w", tableName, err)
	}

	return count, nil
}

func (repo *CockRoachRepository) BulkInsertStocks(ctx context.Context, stocks []*models.FormattedStock, tableName string) error {

	err := repo.createTable(ctx, tableName)
	if err != nil {
		return err
	}

	err = repo.bulkInsertToTable(ctx, tableName, stocks)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CockRoachRepository) BulkUpdateStocks(ctx context.Context, stocks []*models.FormattedStock, originalTable, tempTable string) error {
	err := repo.createTable(ctx, tempTable)
	if err != nil {
		return err
	}

	defer func() error {
		err = repo.dropTable(ctx, tempTable)
		if err != nil {
			return err
		}

		return nil
	}()

	err = repo.bulkInsertToTable(ctx, tempTable, stocks)
	if err != nil {
		return err
	}

	count, err := repo.compareTables(ctx, originalTable, tempTable)
	if err != nil {
		return err
	}

	if count > 0 {
		err = repo.mergeTables(ctx, originalTable, tempTable)
		if err != nil {
			return err
		}

		err = repo.updateTable(ctx, originalTable, tempTable)
		if err != nil {
			return err
		}

	}

	deleteCount, err := repo.compareTables(ctx, tempTable, originalTable)
	if err != nil {
		return err
	}

	if deleteCount > 0 {
		err = repo.deleteObsoleteRows(ctx, originalTable, tempTable)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *CockRoachRepository) compareTables(ctx context.Context, originalTable, tempTable string) (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM (SELECT * FROM %s EXCEPT SELECT * FROM %s)`, pq.QuoteIdentifier(tempTable), pq.QuoteIdentifier(originalTable))
	err := repo.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to compare tables %s and %s: %w", originalTable, tempTable, err)
	}
	return count, nil
}

func (repo *CockRoachRepository) deleteObsoleteRows(ctx context.Context, originalTable, tempTable string) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		mergeQuery := fmt.Sprintf(`
		DELETE FROM %s 
		WHERE ticker IN (
    		SELECT o.ticker 
    		FROM %s o 
    		LEFT JOIN %s t ON o.ticker = t.ticker 
    		WHERE t.ticker IS NULL
			)`, pq.QuoteIdentifier(originalTable), pq.QuoteIdentifier(originalTable), pq.QuoteIdentifier(tempTable))

		result, err := tx.ExecContext(ctx, mergeQuery)
		if err != nil {
			return fmt.Errorf("error deleting rows in table %s: %w", originalTable, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting rows affected in %s: %w", originalTable, err)
		}

		log.Printf("Deleted %d obsolete stocks into %s table", rowsAffected, originalTable)

		return nil
	})
}

func (repo *CockRoachRepository) mergeTables(ctx context.Context, originalTable, tempTable string) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		mergeQuery := fmt.Sprintf(`
        INSERT INTO %s (ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time)
        SELECT t.ticker, t.target_from, t.target_to, t.company, t.action, t.brokerage, t.rating_from, t.rating_to, t.time
        FROM %s t
        LEFT JOIN stocks s ON t.ticker = s.ticker
        WHERE s.ticker IS NULL`, pq.QuoteIdentifier(originalTable), pq.QuoteIdentifier(tempTable))

		result, err := tx.ExecContext(ctx, mergeQuery)
		if err != nil {
			return fmt.Errorf("error merging in table %s: %w", originalTable, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting rows affected in %s: %w", originalTable, err)
		}

		log.Printf("Merged %d new stocks into %s table", rowsAffected, originalTable)

		return nil
	})
}

func (repo *CockRoachRepository) updateTable(ctx context.Context, originalTable, tempTable string) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		updateQuery := fmt.Sprintf(`
		UPDATE %s s
		SET
    		target_from = t.target_from,
    		target_to = t.target_to,
    		company = t.company,
    		time = t.time,
    		action = t.action,
    		brokerage = t.brokerage,
    		rating_from = t.rating_from,
    		rating_to = t.rating_to
		FROM %s t
		WHERE s.ticker = t.ticker
  			AND (s.target_from != t.target_from OR
       			s.target_to != t.target_to OR
       			s.time != t.time OR
       			s.company != t.company OR
       			s.action != t.action OR
       			s.brokerage != t.brokerage OR
       			s.rating_from != t.rating_from OR
       			s.rating_to != t.rating_to
				);
    	`, pq.QuoteIdentifier(originalTable), pq.QuoteIdentifier(tempTable))

		result, err := tx.ExecContext(ctx, updateQuery)
		if err != nil {
			return fmt.Errorf("error updating table %s: %w", originalTable, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting rows affected in %s: %w", originalTable, err)
		}

		log.Printf("Updated %d stocks in %s table", rowsAffected, originalTable)

		return nil
	})
}

func (repo *CockRoachRepository) bulkInsertToTable(ctx context.Context, tableName string, stocks []*models.FormattedStock) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(pq.CopyIn(tableName, "ticker", "target_from", "target_to",
			"company", "action", "brokerage", "rating_from", "rating_to", "time"))
		if err != nil {
			return fmt.Errorf("error preparing bulk insert statement: %w", err)
		}

		defer stmt.Close()

		for _, stock := range stocks {
			_, err = stmt.ExecContext(ctx, stock.Ticker, stock.TargetFrom, stock.TargetTo, stock.Company,
				stock.Action, stock.Brokerage, stock.RatingFrom, stock.RatingTo, stock.Time)
			if err != nil {
				return fmt.Errorf("error adding item %s to bulk insert: %w", stock.Ticker, err)
			}
		}

		_, err = stmt.ExecContext(ctx)
		if err != nil {
			return fmt.Errorf("error finalizing bulk insert to %s: %w", tableName, err)
		}

		log.Printf("Inserted %d stocks into %s", len(stocks), tableName)

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

		if _, err := tx.ExecContext(ctx, createTableQuery); err != nil {
			return fmt.Errorf("error creating table %s: %w", tableName, err)
		}

		log.Printf("Table %s created successfully", tableName)

		return nil
	})
}

func (repo *CockRoachRepository) dropTable(ctx context.Context, tableName string) error {
	return repo.execInTransaction(ctx, func(tx *sql.Tx) error {
		dropTableQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s", pq.QuoteIdentifier(tableName))

		if _, err := tx.ExecContext(ctx, dropTableQuery); err != nil {
			return fmt.Errorf("error droping table %s: %w", tableName, err)
		}

		log.Printf("Table %s dropped successfully", tableName)

		return nil
	})
}

func filterQueryParams(field, order, search, tableName string, page, limit int) (queryBuilder, error) {
	page, limit = normalizePaginationParams(page, limit)
	offset := (page - 1) * limit

	query, err := generatePaginationQuery(field, order, search, tableName)
	if err != nil {
		return queryBuilder{}, err
	}

	params := make([]any, 0)

	searchParams := getSearchParams(search)
	params = append(params, searchParams...)
	params = append(params, limit, offset)

	return queryBuilder{
		query:  query,
		params: params,
	}, nil
}

func generatePaginationQuery(field, order, search, tableName string) (string, error) {

	baseQuery := fmt.Sprintf(`SELECT * FROM %s`, pq.QuoteIdentifier(tableName))

	orderStm, err := buildOrderStatement(field, order)
	if err != nil {
		return "", err
	}

	searchStm := buildSearchStatement(search)

	query := baseQuery
	if searchStm != "" {
		query += " " + searchStm
	}

	query += " " + orderStm

	var pagination string
	if searchStm != "" {
		pagination = " LIMIT $2 OFFSET $3"
	} else {
		pagination = " LIMIT $1 OFFSET $2"
	}

	return query + pagination, nil
}

// validate max pages
func normalizePaginationParams(page, limit int) (int, int) {
	if limit <= 0 {
		limit = defaultLimit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	if page <= 0 {
		page = defaultPage
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
	search = strings.TrimSpace(search)
	if search == "" {
		return ""
	}
	return "WHERE ticker ILIKE $1 OR company ILIKE $1 OR brokerage ILIKE $1"
}

func getSearchParams(search string) []any {
	search = strings.TrimSpace(search)
	if search == "" {
		return []any{}
	}
	return []any{"%" + search + "%"}
}

func buildOrderStatement(field, order string) (string, error) {
	field = strings.ToLower(strings.TrimSpace(field))
	order = strings.ToUpper(strings.TrimSpace(order))

	if field != "" && !isValidField(field) {
		return "", fmt.Errorf("error building statement: %w", ErrInvalidField)
	}

	if order != "" && !isValidOrder(order) {
		return "", fmt.Errorf("error building statement: %w", ErrInvalidOrder)
	}

	if field == "" {
		field = defaultField
	}

	if order == "" {
		order = defaultOrder
	}

	return fmt.Sprintf("ORDER BY %s %s", field, order), nil
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
			return nil, fmt.Errorf("error scanning rows: %w", err)
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
