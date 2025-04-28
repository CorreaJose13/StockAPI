package analysis

import (
	"log"
	"strings"

	"slices"

	"github.com/CorreaJose13/StockAPI/models"
)

var (
	topBrokerages = []string{
		"JPMorgan Chase & Co.",
		"Evercore ISI",
		"Bank of America",
		"Morgan Stanley",
		"Barclays",
		"Citigroup",
		"UBS Group",
		"The Goldman Sachs Group",
		"Wells Fargo & Company",
		"Deutsche Bank Aktiengesellschaft",
	}
)

const (
	porcChangeWeight = 25.0 / 100
	absChangeWeight  = 15.0 / 100
	timeWeight       = 15.0 / 100
	brokerageWeight  = 10.0 / 100
	ratingWeight     = 20.0 / 100
	ratingDiffWeight = 10.0 / 100
	actionWeight     = 5 / 100
)

type Analysis struct {
	Stocks []*models.FormattedStock
}

type stockMetrics struct {
	brokerageMap  map[string]int
	maxPercChange float64
	minPercChange float64
	maxAbsChange  float64
	minAbsChange  float64
	oldestTime    int64
	newestTime    int64
}

func NewAnalysis(stocks []*models.FormattedStock) *Analysis {
	return &Analysis{
		Stocks: stocks,
	}
}

func (a *Analysis) Analyze() {

	metrics := a.computeStockMetrics()

	for _, stock := range a.Stocks {
		porcChange := percentageChange(stock.TargetFrom, stock.TargetTo)
		absChange := absoluteChange(stock.TargetFrom, stock.TargetTo)
		timeValue := stock.Time.Unix()
		porcChangeScore := normalizeValue(porcChange, metrics.minPercChange, metrics.maxPercChange)
		absChangeScore := normalizeValue(absChange, metrics.minAbsChange, metrics.maxAbsChange)
		timeScore := normalizeValue(float64(timeValue), float64(metrics.oldestTime), float64(metrics.newestTime))
		brokerageScore := a.brokerageScore(metrics.brokerageMap, stock.Brokerage, a.getStocksCount())
		ratingScore := mapRatingToFloat(stock.RatingTo)
		ratingDiffScore := ratingDifference(stock.RatingFrom, stock.RatingTo)
		actionValue := mapActionToFloat(stock.Action)

		overallScore := (porcChangeScore * porcChangeWeight) +
			(absChangeScore * absChangeWeight) +
			(timeScore * timeWeight) +
			(brokerageScore * brokerageWeight) +
			(ratingScore * ratingWeight) +
			(ratingDiffScore * ratingDiffWeight) +
			(actionValue * actionWeight)
		stock.Score = overallScore
	}

	log.Println("analysis completed")

}

func (a *Analysis) computeStockMetrics() *stockMetrics {
	// Initialize variables
	stock := a.Stocks[0]
	frequency := make(map[string]int)
	percChange := percentageChange(stock.TargetFrom, stock.TargetTo)
	absChange := absoluteChange(stock.TargetFrom, stock.TargetTo)
	timeValue := stock.Time.Unix()

	minPerc, maxPerc := percChange, percChange
	minAbs, maxAbs := absChange, absChange
	minTime, maxTime := timeValue, timeValue

	for i := 1; i < len(a.Stocks); i++ {
		stock = a.Stocks[i]
		percChange = percentageChange(stock.TargetFrom, stock.TargetTo)
		absChange = absoluteChange(stock.TargetFrom, stock.TargetTo)
		timeValue = stock.Time.Unix()

		normalizedBrokerage := strings.TrimSpace(strings.ToLower(stock.Brokerage))
		frequency[normalizedBrokerage]++

		// Min and max percentage change
		if percChange < minPerc {
			minPerc = percChange
		}

		if percChange > maxPerc {
			maxPerc = percChange
		}

		// Min and max absolute change
		if absChange < minAbs {
			minAbs = absChange
		}

		if absChange > maxAbs {
			maxAbs = absChange
		}

		// oldest and newest time
		if timeValue < minTime {
			minTime = timeValue
		}

		if timeValue > maxTime {
			maxTime = timeValue
		}
	}
	return &stockMetrics{
		brokerageMap:  frequency,
		maxPercChange: maxPerc,
		minPercChange: minPerc,
		maxAbsChange:  maxAbs,
		minAbsChange:  minAbs,
		oldestTime:    minTime,
		newestTime:    maxTime,
	}
}

func normalizeValue(value, min, max float64) float64 {
	if min == max {
		return 1.0
	}
	return (value - min) / (max - min)
}

func percentageChange(targetFrom, targetTo float64) float64 {
	targetDiff := targetTo - targetFrom
	return (targetDiff / targetFrom) * 100
}

func absoluteChange(targetFrom, targetTo float64) float64 {
	return targetTo - targetFrom
}

func isTopBrokerage(brokerage string) bool {
	return slices.Contains(topBrokerages, brokerage)
}

func brokerageRating(brokerage string) float64 {
	if isTopBrokerage(brokerage) {
		return 1.0
	}
	return 0.75
}

func (a *Analysis) brokerageRelativeFrequency(brokerageFrequencyMap map[string]int, brokerage string) float64 {
	normalizedBrokerage := strings.TrimSpace(strings.ToLower(brokerage))
	brokFreq := brokerageFrequencyMap[normalizedBrokerage]
	return float64(brokFreq) / float64(a.getStocksCount())
}

func (a *Analysis) brokerageScore(brokerageFrequencyMap map[string]int, brokerage string, lenStocks int) float64 {
	bRating := brokerageRating(brokerage)
	bRelFreq := a.brokerageRelativeFrequency(brokerageFrequencyMap, brokerage)
	return bRating * bRelFreq
}

func mapRatingToFloat(rating string) float64 {
	switch rating {
	case "buy", "strong-buy", "speculative buy":
		return 1

	case "outperform", "sector outperform", "market outperform", "overweight", "outperformer", "positive":
		return 0.75

	case "hold", "neutral", "unchanged", "market perform", "equal weight", "in-line", "sector perform",
		"sector weight", "peer perform":
		return 0.5

	case "underperform", "underweight", "sector underperform", "under perform", "reduce":
		return 0.25

	case "sell", "negative":
		return 0

	default:
		return 0
	}
}

func ratingDifference(ratingFrom, ratingTo string) float64 {
	ratingFromValue := mapRatingToFloat(ratingFrom)
	ratingToValue := mapRatingToFloat(ratingTo)

	ratingDiff := ratingToValue - ratingFromValue
	if ratingDiff < 0 {
		return 0
	}

	return ratingDiff
}

func mapActionToFloat(action string) float64 {
	switch action {
	case "upgraded by":
		return 1

	case "target raised by":
		return 0.75

	case "initiated by", "reiterated by", "target set by":
		return 0.5

	case "target lowered by":
		return 0.25

	case "downgraded by":
		return 0

	default:
		return 0
	}
}
