package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/CorreaJose13/StockAPI/models"
)

func Formatter(stock *models.Stock) (*models.FormattedStock, error) {
	formattedTargetFrom, err := formatTarget(stock.TargetFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to format target from '%s': %w", stock.TargetFrom, err)
	}

	formattedTargetTo, err := formatTarget(stock.TargetTo)
	if err != nil {
		return nil, fmt.Errorf("failed to format target to '%s': %w", stock.TargetTo, err)
	}

	formattedTime, err := formatTime(stock.Time)
	if err != nil {
		return nil, fmt.Errorf("failed to format time '%s': %w", stock.Time, err)
	}

	return &models.FormattedStock{
		Ticker:     stock.Ticker,
		TargetFrom: formattedTargetFrom,
		TargetTo:   formattedTargetTo,
		Company:    stock.Company,
		Action:     stock.Action,
		Brokerage:  stock.Brokerage,
		RatingFrom: stock.RatingFrom,
		RatingTo:   stock.RatingTo,
		Time:       formattedTime,
	}, nil
}

func formatTarget(target string) (float64, error) {
	targetWithoutDollar := strings.ReplaceAll(target, "$", "")
	targetWithoutCommas := strings.ReplaceAll(targetWithoutDollar, ",", "")

	value, err := strconv.ParseFloat(targetWithoutCommas, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse target price '%s': %w", targetWithoutDollar, err)
	}

	return value, nil
}

func formatTime(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp '%s': %w", timeStr, err)
	}

	return parsedTime, nil
}
