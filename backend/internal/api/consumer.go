package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/models"
)

type ApiConsumer struct {
	client    *http.Client
	apiUrl    string
	authToken string
}

func NewApiConsumer(cfg *config.Config) *ApiConsumer {
	return &ApiConsumer{
		client:    &http.Client{},
		apiUrl:    cfg.ApiUrl,
		authToken: "Bearer " + cfg.BearerToken,
	}
}

func (ac *ApiConsumer) FetchStocks() ([]models.Stock, error) {
	var stocks []models.Stock
	nextPage := ""

	for {
		url := ac.apiUrl
		if nextPage != "" {
			url += "?next_page=" + nextPage
		}

		log.Printf("Fetching stocks from: %s", url)

		body, err := ac.doRequest(url)
		if err != nil {
			return nil, fmt.Errorf("error fetching stocks: %w", err)
		}

		stocks = append(stocks, body.Items...)

		if body.Next_page == "" {
			log.Println("All stock data retrieved successfully")
			break
		}

		nextPage = body.Next_page
	}

	return stocks, nil
}

func (ac *ApiConsumer) doRequest(url string) (*models.Response, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ac.authToken)

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
		return nil, fmt.Errorf("API returned status code: %d and body: %s", resp.StatusCode, string(body))
	}

	var response models.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &response, nil
}
