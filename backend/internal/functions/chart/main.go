package main

import (
	"context"
	"errors"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/internal/api/response"
	"github.com/CorreaJose13/StockAPI/internal/chart"
	"github.com/CorreaJose13/StockAPI/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	ErrMissingAPIKey = errors.New("api key cannot be empty")
)

const (
	maxResults = 10
)

type chartResponse struct {
	TimeSeries []models.DailyData `json:"time_series"`
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg := config.LoadAPIConfig()
	if cfg.APIKEY == "" {
		return events.APIGatewayProxyResponse{}, ErrMissingAPIKey
	}

	chart := chart.NewChartConsumer(cfg)

	ticker := req.QueryStringParameters["ticker"]

	stockData, err := chart.FetchData(ticker)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	chartResponse := chartResponse{
		TimeSeries: stockData[len(stockData)-maxResults-1 : len(stockData)-1],
	}

	return response.Success(chartResponse)
}

func main() {
	lambda.Start(handler)
}
