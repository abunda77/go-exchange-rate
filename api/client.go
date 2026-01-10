package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const BaseURL = "https://use.api.co.id"

// Client represents the API client for exchange rates
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// RateEntry represents a single currency's rates
type RateEntry struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

// AllRatesResponse represents the response from /api/exchange-rates
type AllRatesResponse struct {
	Success           bool        `json:"success"`
	UpdatedAt         int64       `json:"updated_at"`
	Rates             []RateEntry `json:"rates"`
	LastDataUpdatedAt int64       `json:"last_data_updated_at"`
	Message           string      `json:"message,omitempty"`
}

// CurrencyResponse represents the response from /currency/:currency
type CurrencyResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message,omitempty"`
}

// PairResponse represents the response from /currency/exchange-rate
type PairResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message,omitempty"`
}

// NewClient creates a new API client
func NewClient() *Client {
	apiKey := os.Getenv("API_KEY")
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
	}
}

// doRequest performs an HTTP request with the API key header
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("x-api-co-id", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetAllExchangeRates fetches all exchange rates
func (c *Client) GetAllExchangeRates() (*AllRatesResponse, error) {
	url := fmt.Sprintf("%s/api/exchange-rates", BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response AllRatesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetCurrencyRate fetches exchange rate for a specific base currency
func (c *Client) GetCurrencyRate(currency string) (*CurrencyResponse, error) {
	url := fmt.Sprintf("%s/currency/%s", BaseURL, currency)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response CurrencyResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetExchangeRatePairs fetches exchange rate for specific currency pairs
func (c *Client) GetExchangeRatePairs(pairs string) (*PairResponse, error) {
	url := fmt.Sprintf("%s/currency/exchange-rate?pair=%s", BaseURL, pairs)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response PairResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
