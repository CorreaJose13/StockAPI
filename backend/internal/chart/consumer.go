package chart

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/models"
)

type chartConsumer struct {
	client *http.Client
	apiKey string
}

func NewChartConsumer(cfg *config.Config) *chartConsumer {
	return &chartConsumer{
		client: &http.Client{},
		apiKey: cfg.APIKEY,
	}
}

func (ac *chartConsumer) FetchData(ticker string) ([]models.DailyData, error) {
	var timeSeries []models.DailyData

	stockData, err := ac.doRequest(ticker)
	if err != nil {
		return nil, fmt.Errorf("error fetching stocks: %w", err)
	}

	for date, data := range stockData.TimeSeriesDaily {
		open, _ := strconv.ParseFloat(data.Open, 64)
		high, _ := strconv.ParseFloat(data.High, 64)
		low, _ := strconv.ParseFloat(data.Low, 64)
		close, _ := strconv.ParseFloat(data.Close, 64)
		volume, _ := strconv.ParseInt(data.Volume, 10, 64)

		dailyData := models.DailyData{
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
			Date:   date,
		}

		timeSeries = append(timeSeries, dailyData)
	}

	sort.Slice(timeSeries, func(i, j int) bool {
		return timeSeries[i].Date < timeSeries[j].Date
	})

	return timeSeries, nil
}

func (ac *chartConsumer) doRequest(ticker string) (*models.StockData, error) {

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", ticker, ac.apiKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := ac.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status code: %d and body: %s", resp.StatusCode, string(body))
	}

	var stockData models.StockData

	if err := json.Unmarshal(body, &stockData); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &stockData, nil
}
