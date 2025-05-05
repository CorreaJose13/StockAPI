package utils

import (
	"testing"
	"time"

	"github.com/CorreaJose13/StockAPI/models"
	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	tests := []struct {
		name     string
		input    *models.Stock
		expected *models.FormattedStock
		wantErr  bool
	}{
		{
			name: "Valid stock data",
			input: &models.Stock{
				Ticker:     "AAPL",
				TargetFrom: "$150.00",
				TargetTo:   "$170.50",
				Company:    "Apple Inc.",
				Action:     "Upgraded",
				Brokerage:  "Morgan Stanley",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				Time:       "2023-05-10T15:04:05Z",
			},
			expected: &models.FormattedStock{
				Ticker:     "AAPL",
				TargetFrom: 150.00,
				TargetTo:   170.50,
				Company:    "Apple Inc.",
				Action:     "upgraded",
				Brokerage:  "Morgan Stanley",
				RatingFrom: "hold",
				RatingTo:   "buy",
				Time:       time.Date(2023, time.May, 10, 15, 4, 5, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "Stock with commas in target prices",
			input: &models.Stock{
				Ticker:     "AMZN",
				TargetFrom: "$1,200.00",
				TargetTo:   "$1,500.00",
				Company:    "Amazon.com Inc.",
				Action:     "Downgraded",
				Brokerage:  "Goldman Sachs",
				RatingFrom: "Buy",
				RatingTo:   "Neutral",
				Time:       "2023-06-15T10:30:00Z",
			},
			expected: &models.FormattedStock{
				Ticker:     "AMZN",
				TargetFrom: 1200.00,
				TargetTo:   1500.00,
				Company:    "Amazon.com Inc.",
				Action:     "downgraded",
				Brokerage:  "Goldman Sachs",
				RatingFrom: "buy",
				RatingTo:   "neutral",
				Time:       time.Date(2023, time.June, 15, 10, 30, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "Empty action and ratings",
			input: &models.Stock{
				Ticker:     "MSFT",
				TargetFrom: "$300",
				TargetTo:   "$350",
				Company:    "Microsoft Corporation",
				Action:     "",
				Brokerage:  "JP Morgan",
				RatingFrom: "",
				RatingTo:   "",
				Time:       "2023-07-20T09:15:30Z",
			},
			expected: &models.FormattedStock{
				Ticker:     "MSFT",
				TargetFrom: 300.00,
				TargetTo:   350.00,
				Company:    "Microsoft Corporation",
				Action:     "initiated by",
				Brokerage:  "JP Morgan",
				RatingFrom: "neutral",
				RatingTo:   "neutral",
				Time:       time.Date(2023, time.July, 20, 9, 15, 30, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "Invalid target from",
			input: &models.Stock{
				Ticker:     "GOOG",
				TargetFrom: "invalid",
				TargetTo:   "$2000",
				Company:    "Alphabet Inc.",
				Action:     "Maintained",
				Brokerage:  "UBS",
				RatingFrom: "Buy",
				RatingTo:   "Buy",
				Time:       "2023-08-05T14:20:15Z",
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Invalid target to",
			input: &models.Stock{
				Ticker:     "TSLA",
				TargetFrom: "$800",
				TargetTo:   "not-a-number",
				Company:    "Tesla, Inc.",
				Action:     "Upgraded",
				Brokerage:  "Citigroup",
				RatingFrom: "Sell",
				RatingTo:   "Hold",
				Time:       "2023-09-12T11:45:30Z",
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Invalid time format",
			input: &models.Stock{
				Ticker:     "FB",
				TargetFrom: "$350",
				TargetTo:   "$400",
				Company:    "Meta Platforms, Inc.",
				Action:     "Upgraded",
				Brokerage:  "Barclays",
				RatingFrom: "Hold",
				RatingTo:   "Overweight",
				Time:       "invalid-time",
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Formatter(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Ticker, result.Ticker)
			assert.Equal(t, tt.expected.TargetFrom, result.TargetFrom)
			assert.Equal(t, tt.expected.TargetTo, result.TargetTo)
			assert.Equal(t, tt.expected.Company, result.Company)
			assert.Equal(t, tt.expected.Action, result.Action)
			assert.Equal(t, tt.expected.Brokerage, result.Brokerage)
			assert.Equal(t, tt.expected.RatingFrom, result.RatingFrom)
			assert.Equal(t, tt.expected.RatingTo, result.RatingTo)
			assert.Equal(t, tt.expected.Time, result.Time)
		})
	}
}

func TestFormatTarget(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "Simple dollar amount",
			input:    "$100",
			expected: 100.0,
			wantErr:  false,
		},
		{
			name:     "Dollar amount with decimal",
			input:    "$199.99",
			expected: 199.99,
			wantErr:  false,
		},
		{
			name:     "Dollar amount with comma",
			input:    "$1,000.50",
			expected: 1000.50,
			wantErr:  false,
		},
		{
			name:     "Multiple commas",
			input:    "$1,234,567.89",
			expected: 1234567.89,
			wantErr:  false,
		},
		{
			name:     "No dollar sign",
			input:    "500",
			expected: 500.0,
			wantErr:  false,
		},
		{
			name:    "Invalid number",
			input:   "not-a-price",
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatTarget(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatRating(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Lowercase rating",
			input:    "buy",
			expected: "buy",
		},
		{
			name:     "Uppercase rating",
			input:    "SELL",
			expected: "sell",
		},
		{
			name:     "Mixed case rating",
			input:    "OuTpErForM",
			expected: "outperform",
		},
		{
			name:     "Rating with spaces",
			input:    "  hold  ",
			expected: "hold",
		},
		{
			name:     "Empty rating",
			input:    "",
			expected: "neutral",
		},
		{
			name:     "Only spaces",
			input:    "   ",
			expected: "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatRating(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatAction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Lowercase action",
			input:    "upgraded",
			expected: "upgraded",
		},
		{
			name:     "Uppercase action",
			input:    "DOWNGRADED",
			expected: "downgraded",
		},
		{
			name:     "Mixed case action",
			input:    "MaInTaInEd",
			expected: "maintained",
		},
		{
			name:     "Action with spaces",
			input:    "  initiated by  ",
			expected: "initiated by",
		},
		{
			name:     "Empty action",
			input:    "",
			expected: "initiated by",
		},
		{
			name:     "Only spaces",
			input:    "   ",
			expected: "initiated by",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatAction(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "Valid RFC3339 time",
			input:    "2023-10-15T14:30:45Z",
			expected: time.Date(2023, time.October, 15, 14, 30, 45, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Valid RFC3339Nano time",
			input:    "2023-11-20T09:15:30.123456789Z",
			expected: time.Date(2023, time.November, 20, 9, 15, 30, 123456789, time.UTC),
			wantErr:  false,
		},
		{
			name:    "Invalid time format",
			input:   "2023/12/25 12:00:00",
			wantErr: true,
		},
		{
			name:    "Empty time string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatTime(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Format(time.RFC3339Nano), result.Format(time.RFC3339Nano))
		})
	}
}
