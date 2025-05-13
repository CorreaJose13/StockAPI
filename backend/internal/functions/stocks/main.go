package main

import (
	"context"
	"errors"
	"fmt"
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

	ErrInvalidType = errors.New("invalid type")
)

func init() {
	repo, initErr = functions.DBSetup()
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if initErr != nil {
		return response.Error(http.StatusInternalServerError, initErr.Error())
	}

	page, err := strconv.Atoi(req.QueryStringParameters["page"])
	if err != nil {
		return response.Error(http.StatusBadRequest, fmt.Sprintf("%v: page must be a number", ErrInvalidType))
	}
	limit, err := strconv.Atoi(req.QueryStringParameters["limit"])
	if err != nil {
		return response.Error(http.StatusBadRequest, fmt.Sprintf("%v: limit must be a number", ErrInvalidType))
	}

	field := req.QueryStringParameters["field"]
	order := req.QueryStringParameters["order"]
	search := req.QueryStringParameters["search"]

	stocks, err := repository.GetStocksFiltered(ctx, field, order, search, "stocks", page, limit)
	if err != nil {
		if errors.Is(err, db.ErrInvalidField) || errors.Is(err, db.ErrInvalidOrder) {
			return response.Error(http.StatusBadRequest, err.Error())
		}
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	stocksLength, err := repository.GetTableLength(ctx, "stocks")
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	responseBody := map[string]any{
		"stocks": stocks,
		"length": stocksLength,
	}

	return response.Success(responseBody)
}

func main() {
	lambda.Start(handler)
}
