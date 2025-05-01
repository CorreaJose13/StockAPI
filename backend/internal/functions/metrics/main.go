package main

import (
	"context"
	"net/http"

	"github.com/CorreaJose13/StockAPI/internal/analysis"
	"github.com/CorreaJose13/StockAPI/internal/api/response"
	"github.com/CorreaJose13/StockAPI/internal/db"
	"github.com/CorreaJose13/StockAPI/internal/functions"
	"github.com/CorreaJose13/StockAPI/internal/repository"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	repo    *db.CockRoachRepository
	initErr error
)

func init() {
	repo, initErr = functions.Setup()
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if initErr != nil {
		return events.APIGatewayProxyResponse{}, initErr
	}

	defer repo.Close()

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
