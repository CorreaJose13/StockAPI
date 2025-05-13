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

type apiConsumer struct {
	client    *http.Client
	apiURL    string
	authToken string
}

func NewAPIConsumer(cfg *config.Config) *apiConsumer {
	return &apiConsumer{
		client:    &http.Client{},
		apiURL:    cfg.APIURL,
		authToken: "Bearer " + cfg.BearerToken,
	}
}

func (ac *apiConsumer) FetchStocks() ([]models.Stock, error) {
	var stocks []models.Stock
	nextPage := ""

	for {
		url := ac.apiURL
		if nextPage != "" {
			url += "?next_page=" + nextPage
		}

		body, err := ac.doRequest(url)
		if err != nil {
			return nil, fmt.Errorf("error fetching stocks: %w", err)
		}

		log.Printf("Fetched %d stocks from page %s", len(body.Items), nextPage)

		stocks = append(stocks, body.Items...)

		if body.NextPage == "" {
			break
		}

		nextPage = body.NextPage
	}

	return stocks, nil
}

func (ac *apiConsumer) doRequest(url string) (*models.Response, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

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
		return nil, fmt.Errorf("api returned status code: %d and body: %s", resp.StatusCode, string(body))
	}

	var response models.Response

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &response, nil
}
