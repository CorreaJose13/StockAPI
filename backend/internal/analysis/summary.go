package analysis

type StockSummary struct {
	TotalStocks    int `json:"total_stocks"`
	PositiveChange int `json:"positive_change"`
	NegativeChange int `json:"negative_change"`
	NoChange       int `json:"no_change"`
}

func (a *Analysis) GetSummary() *StockSummary {
	countPositive, countNegative, countNeutral := a.countStocksByChangeTrend()
	return &StockSummary{
		TotalStocks:    a.getStocksCount(),
		PositiveChange: countPositive,
		NegativeChange: countNegative,
		NoChange:       countNeutral,
	}
}

func (a *Analysis) getStocksCount() int {
	return len(a.Stocks)
}

func (a *Analysis) countStocksByChangeTrend() (countPositive, countNegative, countNeutral int) {
	countPositive = 0
	countNegative = 0
	countNeutral = 0
	for _, stock := range a.Stocks {
		change := percentageChange(stock.TargetFrom, stock.TargetTo)
		if change > 0 {
			countPositive++
		}
		if change < 0 {
			countNegative++
		}
		if change == 0 {
			countNeutral++
		}
	}
	return
}
