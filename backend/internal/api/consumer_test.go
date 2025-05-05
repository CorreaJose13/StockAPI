package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/CorreaJose13/StockAPI/models"
)

func TestNewAPIConsumer(t *testing.T) {
	cfg := &config.Config{
		APIURL:      "https://api.example.com",
		BearerToken: "test-token",
	}

	consumer := NewAPIConsumer(cfg)

	if consumer.apiURL != cfg.APIURL {
		t.Errorf("Expected apiURL to be %s, got %s", cfg.APIURL, consumer.apiURL)
	}

	expectedToken := "Bearer test-token"
	if consumer.authToken != expectedToken {
		t.Errorf("Expected authToken to be %s, got %s", expectedToken, consumer.authToken)
	}

	if consumer.client == nil {
		t.Error("Expected http client to be initialized")
	}
}

func TestDoRequest(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
		}

		response := models.Response{
			Items: []models.Stock{
				{
					Ticker:     "AAPL",
					TargetFrom: "$150",
					TargetTo:   "$170",
					Company:    "Apple Inc.",
				},
			},
			NextPage: "",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		BearerToken: "test-token",
	}
	consumer := NewAPIConsumer(cfg)

	response, err := consumer.doRequest(server.URL)
	if err != nil {
		t.Fatalf("doRequest returned unexpected error: %v", err)
	}

	if len(response.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(response.Items))
	}

}

func TestFetchStocks(t *testing.T) {

	pageCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var nextPage string
		var items []models.Stock

		switch pageCount {
		case 0:
			// Primera página
			nextPage = "page2"
			items = []models.Stock{
				{
					Ticker:     "AAPL",
					TargetFrom: "$150",
					TargetTo:   "$170",
					Company:    "Apple Inc.",
					Action:     "upgraded by",
					Brokerage:  "Morgan Stanley",
					RatingFrom: "hold",
					RatingTo:   "buy",
					Time:       time.Now().Format(time.RFC3339),
				},
			}
		case 1:
			// Segunda página
			nextPage = "page3"
			items = []models.Stock{
				{
					Ticker:     "MSFT",
					TargetFrom: "$280",
					TargetTo:   "$310",
					Company:    "Microsoft Corp",
					Action:     "maintained by",
					Brokerage:  "Goldman Sachs",
					RatingFrom: "buy",
					RatingTo:   "buy",
					Time:       time.Now().Format(time.RFC3339),
				},
			}
		case 2:
			// Tercera y última página
			nextPage = ""
			items = []models.Stock{
				{
					Ticker:     "GOOG",
					TargetFrom: "$2500",
					TargetTo:   "$2700",
					Company:    "Alphabet Inc",
					Action:     "downgraded by",
					Brokerage:  "JP Morgan",
					RatingFrom: "buy",
					RatingTo:   "hold",
					Time:       time.Now().Format(time.RFC3339),
				},
			}
		}

		response := models.Response{
			Items:    items,
			NextPage: nextPage,
		}

		pageCount++
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		BearerToken: "test-token",
	}
	consumer := NewAPIConsumer(cfg)

	stocks, err := consumer.FetchStocks()
	if err != nil {
		t.Fatalf("FetchStocks returned unexpected error: %v", err)
	}

	if len(stocks) != 3 {
		t.Errorf("Expected 3 stocks (from 3 pages), got %d", len(stocks))
	}

}

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Unauthorized"}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		BearerToken: "invalid-token",
	}
	consumer := NewAPIConsumer(cfg)

	_, err := consumer.FetchStocks()
	if err == nil {
		t.Error("Expected error for unauthorized request, got nil")
	}
}

func TestInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"invalid json`))
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		BearerToken: "test-token",
	}
	consumer := NewAPIConsumer(cfg)

	_, err := consumer.FetchStocks()
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}
