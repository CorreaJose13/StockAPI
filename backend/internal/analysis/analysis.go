package analysis

import (
	"sort"
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
	percChangeWeight = 25.0 / 100
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

type StockAnalysis struct {
	*models.FormattedStock
	Score float64 `json:"score"`
}

type StockMetrics struct {
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

func (a *Analysis) Analyze() []*StockAnalysis {

	metrics := a.computeStockMetrics()

	var stocksAnalysis []*StockAnalysis
	for _, stock := range a.Stocks {

		score := a.calculateScore(stock, metrics)

		stocksAnalysis = append(stocksAnalysis, &StockAnalysis{
			FormattedStock: stock,
			Score:          score})
	}

	sort.Slice(stocksAnalysis, func(i, j int) bool {
		return stocksAnalysis[i].Score > stocksAnalysis[j].Score
	})

	resultLimit := min(len(stocksAnalysis), 10)

	return stocksAnalysis[:resultLimit]
}

func (a *Analysis) calculateScore(stock *models.FormattedStock, metrics *StockMetrics) float64 {
	percChange := percentageChange(stock.TargetFrom, stock.TargetTo)
	absChange := absoluteChange(stock.TargetFrom, stock.TargetTo)
	timeValue := stock.Time.Unix()
	percChangeScore := normalizeValue(percChange, metrics.minPercChange, metrics.maxPercChange)
	absChangeScore := normalizeValue(absChange, metrics.minAbsChange, metrics.maxAbsChange)
	timeScore := normalizeValue(float64(timeValue), float64(metrics.oldestTime), float64(metrics.newestTime))
	brokerageScore := a.brokerageScore(metrics.brokerageMap, stock.Brokerage)
	ratingScore := mapRatingToFloat(stock.RatingTo)
	ratingDiffScore := ratingDifference(stock.RatingFrom, stock.RatingTo)
	actionValue := mapActionToFloat(stock.Action)

	overallScore := (percChangeScore * percChangeWeight) +
		(absChangeScore * absChangeWeight) +
		(timeScore * timeWeight) +
		(brokerageScore * brokerageWeight) +
		(ratingScore * ratingWeight) +
		(ratingDiffScore * ratingDiffWeight) +
		(actionValue * actionWeight)

	return overallScore
}

func (a *Analysis) computeStockMetrics() *StockMetrics {
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

		minPerc, maxPerc = setMinMax(percChange, minPerc, maxPerc)
		minAbs, maxAbs = setMinMax(absChange, minAbs, maxAbs)
		minTimeF, maxTimeF := setMinMax(float64(timeValue), float64(minTime), float64(maxTime))
		minTime = int64(minTimeF)
		maxTime = int64(maxTimeF)
	}
	return &StockMetrics{
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

func setMinMax(value, min, max float64) (float64, float64) {
	if value < min {
		min = value
	}
	if value > max {
		max = value
	}
	return min, max
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

func (a *Analysis) brokerageScore(brokerageFrequencyMap map[string]int, brokerage string) float64 {
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
