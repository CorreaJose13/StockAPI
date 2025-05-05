package analysis

import (
	"testing"

	"github.com/CorreaJose13/StockAPI/models"
)

func TestGetSummary(t *testing.T) {
	analysis := &Analysis{
		Stocks: []*models.FormattedStock{
			{Ticker: "AAPL", TargetFrom: 100.0, TargetTo: 110.0}, // Positivo
			{Ticker: "GOOG", TargetFrom: 200.0, TargetTo: 180.0}, // Negativo
			{Ticker: "MSFT", TargetFrom: 150.0, TargetTo: 150.0}, // Sin cambio
			{Ticker: "AMZN", TargetFrom: 300.0, TargetTo: 320.0}, // Positivo
			{Ticker: "META", TargetFrom: 250.0, TargetTo: 225.0}, // Negativo
		},
	}

	summary := analysis.GetSummary()

	if summary.TotalStocks != 5 {
		t.Errorf("Expected total stocks to be 5, got %d", summary.TotalStocks)
	}

	if summary.PositiveChange != 2 {
		t.Errorf("Expected positive changes to be 2, got %d", summary.PositiveChange)
	}

	if summary.NegativeChange != 2 {
		t.Errorf("Expected negative changes to be 2, got %d", summary.NegativeChange)
	}

	if summary.NoChange != 1 {
		t.Errorf("Expected no changes to be 1, got %d", summary.NoChange)
	}
}

func TestGetStocksCount(t *testing.T) {
	tests := []struct {
		name      string
		stocks    []*models.FormattedStock
		wantCount int
	}{
		{
			name:      "Empty stocks",
			stocks:    []*models.FormattedStock{},
			wantCount: 0,
		},
		{
			name: "Single stock",
			stocks: []*models.FormattedStock{
				{Ticker: "AAPL"},
			},
			wantCount: 1,
		},
		{
			name: "Multiple stocks",
			stocks: []*models.FormattedStock{
				{Ticker: "AAPL"},
				{Ticker: "GOOG"},
				{Ticker: "MSFT"},
			},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &Analysis{Stocks: tt.stocks}
			count := analysis.getStocksCount()
			if count != tt.wantCount {
				t.Errorf("getStocksCount() = %v, want %v", count, tt.wantCount)
			}
		})
	}
}

func TestCountStocksByChangeTrend(t *testing.T) {
	tests := []struct {
		name         string
		stocks       []*models.FormattedStock
		wantPositive int
		wantNegative int
		wantNeutral  int
	}{
		{
			name:         "Empty stocks",
			stocks:       []*models.FormattedStock{},
			wantPositive: 0,
			wantNegative: 0,
			wantNeutral:  0,
		},
		{
			name: "Only positive changes",
			stocks: []*models.FormattedStock{
				{TargetFrom: 100, TargetTo: 110},
				{TargetFrom: 200, TargetTo: 250},
			},
			wantPositive: 2,
			wantNegative: 0,
			wantNeutral:  0,
		},
		{
			name: "Only negative changes",
			stocks: []*models.FormattedStock{
				{TargetFrom: 100, TargetTo: 90},
				{TargetFrom: 200, TargetTo: 150},
			},
			wantPositive: 0,
			wantNegative: 2,
			wantNeutral:  0,
		},
		{
			name: "Only neutral changes",
			stocks: []*models.FormattedStock{
				{TargetFrom: 100, TargetTo: 100},
				{TargetFrom: 200, TargetTo: 200},
			},
			wantPositive: 0,
			wantNegative: 0,
			wantNeutral:  2,
		},
		{
			name: "Mixed changes",
			stocks: []*models.FormattedStock{
				{TargetFrom: 100, TargetTo: 120},
				{TargetFrom: 200, TargetTo: 180},
				{TargetFrom: 300, TargetTo: 300},
				{TargetFrom: 400, TargetTo: 450},
			},
			wantPositive: 2,
			wantNegative: 1,
			wantNeutral:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &Analysis{Stocks: tt.stocks}
			positive, negative, neutral := analysis.countStocksByChangeTrend()

			if positive != tt.wantPositive {
				t.Errorf("Positive count = %v, want %v", positive, tt.wantPositive)
			}

			if negative != tt.wantNegative {
				t.Errorf("Negative count = %v, want %v", negative, tt.wantNegative)
			}

			if neutral != tt.wantNeutral {
				t.Errorf("Neutral count = %v, want %v", neutral, tt.wantNeutral)
			}
		})
	}
}
