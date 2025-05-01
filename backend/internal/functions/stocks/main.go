package main

import (
	"context"
	"net/http"
	"strconv"

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
		return response.Error(http.StatusInternalServerError, initErr.Error())
	}

	page, _ := strconv.Atoi(req.QueryStringParameters["page"])
	limit, _ := strconv.Atoi(req.QueryStringParameters["limit"])

	stocks, err := repository.GetStocksPaginated(ctx, page, limit)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(stocks)
}

func main() {
	lambda.Start(handler)
}
