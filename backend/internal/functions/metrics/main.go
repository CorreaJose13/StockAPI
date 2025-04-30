package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/internal/analysis"
	"github.com/CorreaJose13/StockAPI/internal/api/response"
	"github.com/CorreaJose13/StockAPI/internal/db"
	"github.com/CorreaJose13/StockAPI/internal/repository"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	dbString        string
	initErr         error
	errMissingDBURL = errors.New("db url cannot be empty")
)

func init() {
	dbString = os.Getenv("DB_URL")
	if dbString == "" {
		initErr = errMissingDBURL
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if initErr != nil {
		return events.APIGatewayProxyResponse{}, initErr
	}

	cfg := &config.Config{
		DBURL: dbString,
	}

	repo, err := db.NewPostgresRepository(cfg)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer repo.Close()

	repository.SetStockRepository(repo)

	stocks, err := repository.GetStocks(ctx)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	analysis := analysis.NewAnalysis(stocks)

	return response.Success(analysis.GetSummary())
}

func main() {
	lambda.Start(handler)
}
