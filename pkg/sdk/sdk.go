package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type client struct {
	host  string
	token string
	http  *http.Client
}

type NewClientOptions struct {
	Host string
}

// New Client
func NewClient(opts NewClientOptions) *client {
	return &client{
		host: opts.Host,
		http: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Set Auth token
func (c *client) SetAuth(token string) {
	c.token = token
}

// Standardized Get request
func (c *client) get(endpoint string, output interface{}) error {
	req, err := http.NewRequest("GET", c.host+endpoint, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("could not get recipes: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("unexpected response '%d'", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(output)
}

// Standardized Post request
func (c *client) post(endpoint string, input, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("could not marshal input: %w", err)
	}

	req, err := http.NewRequest("POST", c.host+endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("could not post recipes: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("unexpected response '%d', could not read response body: %w", res.StatusCode, err)
		}
		return fmt.Errorf("unexpected response '%d': %s", res.StatusCode, string(body))
	}

	if output != nil {
		return json.NewDecoder(res.Body).Decode(output)
	}

	return nil
}
