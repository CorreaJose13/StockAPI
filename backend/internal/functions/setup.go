package functions

import (
	"errors"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/internal/db"
	"github.com/CorreaJose13/StockAPI/internal/repository"
)

var (
	ErrMissingDBURL       = errors.New("db url cannot be empty")
	ErrMissingAPIURL      = errors.New("api url cannot be empty")
	ErrMissingBearerToken = errors.New("bearer token cannot be empty")
)

func DBSetup() (*db.CockRoachRepository, error) {
	cfg := config.LoadDbConfig()
	if cfg.DBURL == "" {
		return nil, ErrMissingDBURL
	}

	repo, err := db.NewPostgresRepository(cfg)
	if err != nil {
		return nil, err
	}

	repository.SetStockRepository(repo)

	return repo, nil
}

func FullSetup() (*db.CockRoachRepository, *config.Config, error) {
	cfg := config.LoadLambdaConfig()
	if cfg.DBURL == "" {
		return nil, nil, ErrMissingDBURL
	}
	if cfg.APIURL == "" {
		return nil, nil, ErrMissingAPIURL
	}
	if cfg.BearerToken == "" {
		return nil, nil, ErrMissingBearerToken
	}

	repo, err := db.NewPostgresRepository(cfg)
	if err != nil {
		return nil, nil, err
	}

	repository.SetStockRepository(repo)

	return repo, cfg, nil
}
