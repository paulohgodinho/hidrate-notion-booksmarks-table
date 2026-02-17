package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client handles communication with the webmeatscraper service
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new scraper client
func NewClient(baseURL string) *Client {
	// Default to localhost if empty
	if baseURL == "" {
		baseURL = "http://localhost:7878"
	}

	// Ensure baseURL doesn't end with a slash
	baseURL = strings.TrimRight(baseURL, "/")

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 80 * time.Second, // 80 second timeout for scraping requests
		},
	}
}

// ScrapeResult contains both the parsed content and raw JSON response
type ScrapeResult struct {
	Content *ScrapedContent
	RawJSON string
}

// Scrape sends a URL to the scraper service and returns the scraped content
func (c *Client) Scrape(ctx context.Context, url string) (*ScrapeResult, error) {
	if url == "" {
		return nil, fmt.Errorf("URL is required")
	}

	// Create request body
	reqBody := ScrapeRequest{
		URL: url,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/scrape", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scraper returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var content ScrapedContent
	if err := json.Unmarshal(body, &content); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Return both parsed content and raw JSON
	return &ScrapeResult{
		Content: &content,
		RawJSON: string(body),
	}, nil
}

// Health checks if the scraper service is healthy and reachable
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send health check request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read health check response: %w", err)
	}

	var health HealthResponse
	if err := json.Unmarshal(body, &health); err != nil {
		// If we can't parse the response, but status is 200, assume healthy
		return &HealthResponse{Status: "ok"}, nil
	}

	return &health, nil
}

// Exit sends a GET request to the /exit endpoint to signal the scraper service to shut down
func (c *Client) Exit(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/exit", nil)
	if err != nil {
		return fmt.Errorf("failed to create exit request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send exit request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("exit request returned status %d", resp.StatusCode)
	}

	return nil
}

// BaseURL returns the base URL of the scraper service
func (c *Client) BaseURL() string {
	return c.baseURL
}
