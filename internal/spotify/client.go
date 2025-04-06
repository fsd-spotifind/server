package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	config     *Config
}

func NewClient(config *Config) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		config:     config,
	}
}

func (c *Client) request(ctx context.Context, method, endpoint, token string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.config.BaseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) doJSON(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("spotify API error: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
