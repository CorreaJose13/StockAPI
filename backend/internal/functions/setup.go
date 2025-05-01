package functions

import (
	"errors"
	"os"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/internal/db"
	"github.com/CorreaJose13/StockAPI/internal/repository"
)

var (
	ErrMissingDBURL = errors.New("db url cannot be empty")
)

func Setup() (*db.CockRoachRepository, error) {
	dbString := os.Getenv("DB_URL")
	if dbString == "" {
		return nil, ErrMissingDBURL
	}

	cfg := &config.Config{
		DBURL: dbString,
	}

	repo, err := db.NewPostgresRepository(cfg)
	if err != nil {
		return nil, err
	}

	repository.SetStockRepository(repo)

	return repo, nil
}
