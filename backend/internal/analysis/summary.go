package analysis

type StockSummary struct {
	TotalStocks    int `json:"total_stocks"`
	PositiveChange int `json:"positive_change"`
	NegativeChange int `json:"negative_change"`
	NoChange       int `json:"no_change"`
}

func (a *Analysis) GetSummary() *StockSummary {
	return &StockSummary{
		TotalStocks:    a.getStocksCount(),
		PositiveChange: a.positiveChange(),
		NegativeChange: a.negativeChange(),
		NoChange:       a.noChange(),
	}
}

func (a *Analysis) getStocksCount() int {
	return len(a.Stocks)
}

func (a *Analysis) countStocksByChangeTrend(trend string) int {
	count := 0
	for _, stock := range a.Stocks {
		change := percentageChange(stock.TargetFrom, stock.TargetTo)

		switch trend {
		case "positive":
			if change > 0 {
				count++
			}
		case "negative":
			if change < 0 {
				count++
			}
		case "no-change":
			if change == 0 {
				count++
			}
		}
	}
	return count
}

func (a *Analysis) positiveChange() int {
	return a.countStocksByChangeTrend("positive")
}

func (a *Analysis) negativeChange() int {
	return a.countStocksByChangeTrend("negative")
}

func (a *Analysis) noChange() int {
	return a.countStocksByChangeTrend("no-change")
}
