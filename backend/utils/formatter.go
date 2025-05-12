package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/CorreaJose13/StockAPI/models"
)

var (
	DefaultRating = "hold"
	DefaultAction = "initiated by"

	ErrEmptyTickerString    = fmt.Errorf("empty ticker")
	ErrEmptyCompanyString   = fmt.Errorf("empty company")
	ErrEmptyBrokerageString = fmt.Errorf("empty brokerage")
	ErrEmptyTargetString    = fmt.Errorf("empty target")
	ErrNegativeTarget       = fmt.Errorf("negative target")
	ErrEmptyTimeString      = fmt.Errorf("empty timestamp")
	ErrInvalidTimeFormat    = fmt.Errorf("invalid format")
)

func Formatter(stock *models.Stock) (*models.FormattedStock, error) {

	formattedTicker, err := formatTicker(stock.Ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to format ticker for '%s': %w", stock.Ticker, err)
	}

	formattedCompany, err := formatCompany(stock.Company)
	if err != nil {
		return nil, fmt.Errorf("failed to format company for '%s': %w", stock.Ticker, err)
	}

	formattedBrokerage, err := formatBrokerage(stock.Brokerage)
	if err != nil {
		return nil, fmt.Errorf("failed to format brokerage for '%s': %w", stock.Ticker, err)
	}

	formattedRatingFrom, err := formatRating(stock.RatingFrom)
	if err != nil {
		log.Printf("warning: empty rating_from for '%s': %v", stock.Ticker, err)
	}

	formattedRatingTo, err := formatRating(stock.RatingTo)
	if err != nil {
		log.Printf("warning: empty rating_to for '%s': %v", stock.Ticker, err)
	}

	formattedAction, err := formatAction(stock.Action)
	if err != nil {
		log.Printf("warning: empty action for '%s': %v", stock.Ticker, err)
	}

	formattedTargetFrom, err := formatTarget(stock.TargetFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to format target_from for '%s': %w", stock.Ticker, err)
	}

	formattedTargetTo, err := formatTarget(stock.TargetTo)
	if err != nil {
		return nil, fmt.Errorf("failed to format target_to for '%s': %w", stock.Ticker, err)
	}

	formattedTime, err := formatTime(stock.Time)
	if err != nil {
		return nil, fmt.Errorf("failed to format time for '%s': %w", stock.Ticker, err)
	}

	return &models.FormattedStock{
		Ticker:     formattedTicker,
		TargetFrom: formattedTargetFrom,
		TargetTo:   formattedTargetTo,
		Company:    formattedCompany,
		Action:     formattedAction,
		Brokerage:  formattedBrokerage,
		RatingFrom: formattedRatingFrom,
		RatingTo:   formattedRatingTo,
		Time:       formattedTime,
	}, nil
}

func formatField(field string, fieldValue string, err error) (string, error) {
	field = strings.TrimSpace(field)
	if len(field) == 0 {
		return "", fmt.Errorf("%w: %s name is required", err, field)
	}
	return fieldValue, nil
}

func formatTicker(ticker string) (string, error) {
	fmtTicker, err := formatField("ticker", ticker, ErrEmptyTickerString)
	return strings.ToUpper(fmtTicker), err
}

func formatCompany(company string) (string, error) {
	return formatField("company", company, ErrEmptyCompanyString)
}

func formatBrokerage(brokerage string) (string, error) {
	return formatField("brokerage", brokerage, ErrEmptyBrokerageString)
}

func formatTarget(target string) (float64, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return 0, fmt.Errorf("%w: target value is required", ErrEmptyTargetString)
	}

	targetWithoutDollar := strings.ReplaceAll(target, "$", "")
	targetWithoutCommas := strings.ReplaceAll(targetWithoutDollar, ",", "")

	value, err := strconv.ParseFloat(targetWithoutCommas, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: failed to parse target price '%s': %v", ErrInvalidTimeFormat, targetWithoutCommas, err)
	}

	if value < 0 {
		return 0, fmt.Errorf("%w: target price cannot be negative: %v", ErrNegativeTarget, value)
	}

	return value, nil
}

func formatDefaultField(field string, fieldValue string, defaultValue string) (string, error) {
	field = strings.TrimSpace(strings.ToLower(field))
	if len(field) == 0 {
		return defaultValue, fmt.Errorf("empty %s  defaulted to '%s'", field, defaultValue)
	}

	return fieldValue, nil
}

func formatRating(rating string) (string, error) {
	formattedRating, err := formatDefaultField("rating", rating, DefaultRating)
	return narrowRating(formattedRating), err
}

func formatAction(action string) (string, error) {
	return formatDefaultField("action", action, DefaultAction)
}

func formatTime(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("%w: timestamp is required", ErrEmptyTimeString)
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: failed to parse timestamp '%s': %v", ErrInvalidTimeFormat, timeStr, err)
	}

	return parsedTime, nil
}

func narrowRating(rating string) string {
	switch rating {
	case "buy", "strong-buy", "positive":
		return "buy"

	case "outperform", "sector outperform", "market outperform", "overweight", "outperformer", "speculative buy":
		return "outperform"

	case "hold", "neutral", "unchanged", "market perform", "equal weight", "in-line", "sector perform",
		"sector weight", "peer perform":
		return "hold"

	case "underperform", "underweight", "sector underperform", "under perform", "reduce":
		return "underperform"

	case "sell", "negative":
		return "sell"
	default:
		return rating
	}
}
