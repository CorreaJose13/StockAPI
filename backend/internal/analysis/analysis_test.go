package analysis

import (
	"fmt"
	"testing"
	"time"

	"github.com/CorreaJose13/StockAPI/models"
)

func TestPercentageChange(t *testing.T) {
	tests := []struct {
		name       string
		targetFrom float64
		targetTo   float64
		want       float64
	}{
		{"Positive change", 100, 120, 20},
		{"Negative change", 100, 80, -20},
		{"No change", 100, 100, 0},
		{"Zero from", 0, 50, 50},
		{"Zero to", 100, 0, -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := percentageChange(tt.targetFrom, tt.targetTo)
			if got != tt.want {
				t.Errorf("percentageChange(%v, %v) = %v, want %v", tt.targetFrom, tt.targetTo, got, tt.want)
			}
		})
	}
}

func TestAbsoluteChange(t *testing.T) {
	tests := []struct {
		name       string
		targetFrom float64
		targetTo   float64
		want       float64
	}{
		{"Positive change", 100, 120, 20},
		{"Negative change", 100, 80, -20},
		{"No change", 100, 100, 0},
		{"Zero values", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := absoluteChange(tt.targetFrom, tt.targetTo); got != tt.want {
				t.Errorf("absoluteChange(%v, %v) = %v, want %v", tt.targetFrom, tt.targetTo, got, tt.want)
			}
		})
	}
}

func TestSetMinMax(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		min     float64
		max     float64
		wantMin float64
		wantMax float64
	}{
		{"Value below min", 5, 10, 20, 5, 20},
		{"Value above max", 25, 10, 20, 10, 25},
		{"Value between", 15, 10, 20, 10, 20},
		{"Value equals min", 10, 10, 20, 10, 20},
		{"Value equals max", 20, 10, 20, 10, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := setMinMax(tt.value, tt.min, tt.max)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("setMinMax(%v, %v, %v) = (%v, %v), want (%v, %v)",
					tt.value, tt.min, tt.max, gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestNormalizeValue(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		min   float64
		max   float64
		want  float64
	}{
		{"Min value", 10, 10, 20, 0},
		{"Max value", 20, 10, 20, 1},
		{"Middle value", 15, 10, 20, 0.5},
		{"Min equals max", 15, 15, 15, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeValue(tt.value, tt.min, tt.max); got != tt.want {
				t.Errorf("normalizeValue(%v, %v, %v) = %v, want %v",
					tt.value, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestIsTopBrokerage(t *testing.T) {
	tests := []struct {
		name      string
		brokerage string
		want      bool
	}{
		{"Top brokerage", "JPMorgan Chase & Co.", true},
		{"Not top brokerage", "Small Firm Inc.", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTopBrokerage(tt.brokerage); got != tt.want {
				t.Errorf("isTopBrokerage(%v) = %v, want %v", tt.brokerage, got, tt.want)
			}
		})
	}
}

func TestMapRatingToFloat(t *testing.T) {
	tests := []struct {
		name   string
		rating string
		want   float64
	}{
		{"Buy rating", "buy", 1},
		{"Outperform rating", "outperform", 0.75},
		{"Hold rating", "hold", 0.5},
		{"Underperform rating", "underperform", 0.25},
		{"Sell rating", "sell", 0},
		{"Unknown rating", "unknown", 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapRatingToFloat(tt.rating); got != tt.want {
				t.Errorf("mapRatingToFloat(%v) = %v, want %v", tt.rating, got, tt.want)
			}
		})
	}
}

func TestMapActionToFloat(t *testing.T) {
	tests := []struct {
		name   string
		action string
		want   float64
	}{
		{"Upgraded action", "upgraded by", 1},
		{"Target raised action", "target raised by", 0.75},
		{"Initiated action", "initiated by", 0.5},
		{"Target lowered action", "target lowered by", 0.25},
		{"Downgraded action", "downgraded by", 0},
		{"Unknown action", "unknown", 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapActionToFloat(tt.action); got != tt.want {
				t.Errorf("mapActionToFloat(%v) = %v, want %v", tt.action, got, tt.want)
			}
		})
	}
}

func TestRatingDifference(t *testing.T) {
	tests := []struct {
		name       string
		ratingFrom string
		ratingTo   string
		want       float64
	}{
		{"Positive change", "hold", "buy", 0.5},
		{"Negative change", "buy", "hold", 0},
		{"No change", "buy", "buy", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ratingDifference(tt.ratingFrom, tt.ratingTo); got != tt.want {
				t.Errorf("ratingDifference(%v, %v) = %v, want %v", tt.ratingFrom, tt.ratingTo, got, tt.want)
			}
		})
	}
}

func TestAnalyze(t *testing.T) {

	now := time.Now()
	stocks := []*models.FormattedStock{
		{
			Ticker:     "AAPL",
			TargetFrom: 100,
			TargetTo:   120,
			Company:    "Apple Inc.",
			Action:     "upgraded by",
			Brokerage:  "JPMorgan Chase & Co.",
			RatingFrom: "hold",
			RatingTo:   "buy",
			Time:       now,
		},
		{
			Ticker:     "MSFT",
			TargetFrom: 200,
			TargetTo:   180,
			Company:    "Microsoft Corp.",
			Action:     "downgraded by",
			Brokerage:  "Citigroup",
			RatingFrom: "buy",
			RatingTo:   "hold",
			Time:       now.Add(-24 * time.Hour),
		},
		{
			Ticker:     "GOOG",
			TargetFrom: 150,
			TargetTo:   150,
			Company:    "Alphabet Inc.",
			Action:     "reiterated by",
			Brokerage:  "Small Firm Inc.",
			RatingFrom: "neutral",
			RatingTo:   "neutral",
			Time:       now.Add(-48 * time.Hour),
		},
	}

	analysis := NewAnalysis(stocks)
	results := analysis.Analyze()

	if len(results) == 0 {
		t.Error("Expected non-empty results from Analyze(), got empty slice")
	}

	for i := range len(results) - 1 {
		if results[i].Score < results[i+1].Score {
			t.Errorf("Results not properly sorted. Score at index %d (%f) should be >= score at index %d (%f)",
				i, results[i].Score, i+1, results[i+1].Score)
		}
	}

	if len(stocks) > limitAnalysis && len(results) > limitAnalysis {
		t.Errorf("Expected results to be limited to %d, got %d", limitAnalysis, len(results))
	}
}

func TestBrokerageScore(t *testing.T) {
	stocks := []*models.FormattedStock{
		{Brokerage: "JPMorgan Chase & Co."},
		{Brokerage: "JPMorgan Chase & Co."},
		{Brokerage: "Small Firm Inc."},
	}

	analysis := NewAnalysis(stocks)

	brokerageMap := map[string]int{
		"jpmorgan chase & co.": 2,
		"small firm inc.":      1,
	}

	score1 := analysis.brokerageScore(brokerageMap, "JPMorgan Chase & Co.")
	expectedScore1 := 1.0 * (2.0 / 3.0) // rating * relative frequency

	if score1 != expectedScore1 {
		t.Errorf("brokerageScore for top brokerage = %v, want %v", score1, expectedScore1)
	}

	score2 := analysis.brokerageScore(brokerageMap, "Small Firm Inc.")
	expectedScore2 := 0.75 * (1.0 / 3.0)

	if score2 != expectedScore2 {
		t.Errorf("brokerageScore for non-top brokerage = %v, want %v", score2, expectedScore2)
	}
}

func TestComputeStockMetrics(t *testing.T) {
	now := time.Now()
	stocks := []*models.FormattedStock{
		{
			TargetFrom: 100,
			TargetTo:   120,
			Brokerage:  "JPMorgan Chase & Co.",
			Time:       now,
		},
		{
			TargetFrom: 200,
			TargetTo:   180,
			Brokerage:  "Citigroup",
			Time:       now.Add(-24 * time.Hour),
		},
		{
			TargetFrom: 150,
			TargetTo:   150,
			Brokerage:  "Small Firm Inc.",
			Time:       now.Add(-48 * time.Hour),
		},
	}

	analysis := NewAnalysis(stocks)
	metrics := analysis.computeStockMetrics()

	if metrics.minPercChange != -10.0 { // (180-200)/200 * 100
		t.Errorf("Expected minPercChange = -10.0, got %v", metrics.minPercChange)
	}

	if metrics.maxPercChange != 20.0 { // (120-100)/100 * 100
		t.Errorf("Expected maxPercChange = 20.0, got %v", metrics.maxPercChange)
	}

	if metrics.minAbsChange != -20.0 { // 180-200
		t.Errorf("Expected minAbsChange = -20.0, got %v", metrics.minAbsChange)
	}

	if metrics.maxAbsChange != 20.0 { // 120-100
		t.Errorf("Expected maxAbsChange = 20.0, got %v", metrics.maxAbsChange)
	}

	fmt.Println(metrics.brokerageMap)

	if metrics.brokerageMap["jpmorgan chase & co."] != 1 {
		t.Errorf("Expected brokerageMap['jpmorgan chase & co.'] = 1, got %v",
			metrics.brokerageMap["jpmorgan chase & co."])
	}

	if metrics.brokerageMap["citigroup"] != 1 {
		t.Errorf("Expected brokerageMap['citigroup'] = 1, got %v",
			metrics.brokerageMap["citigroup"])
	}

	if metrics.brokerageMap["small firm inc."] != 1 {
		t.Errorf("Expected brokerageMap['small firm inc.'] = 1, got %v",
			metrics.brokerageMap["small firm inc."])
	}

	oldestTime := now.Add(-48 * time.Hour).Unix()
	if metrics.oldestTime != oldestTime {
		t.Errorf("Expected oldestTime = %v, got %v", oldestTime, metrics.oldestTime)
	}

	newestTime := now.Unix()
	if metrics.newestTime != newestTime {
		t.Errorf("Expected newestTime = %v, got %v", newestTime, metrics.newestTime)
	}
}
