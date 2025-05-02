package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

type SearchClient struct {
	client *http.Client
	token  string
	host   string
}

func NewClient(token, host string) *SearchClient {
	return &SearchClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		token: token,
		host:  host,
	}
}

func (s *SearchClient) Search(ctx context.Context, query string) (SearchResult, error) {
	// Build the query, encode the query, and append the default query parameters
	q := "?q=" + url.QueryEscape(query) + "&safesearch=strict&result_filter=web"

	var result SearchResult
	err := s.get(ctx, "/res/v1/web/search"+q, &result)
	if err != nil {
		return SearchResult{}, err
	}

	log.Info().
		Str("producer", "braveapi_web_search").
		Str("query", query).
		Msg("braveapi metadata")

	return result, nil
}

func (s *SearchClient) get(ctx context.Context, path string, output interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.host+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Subscription-Token", s.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("unexpected response '%d', could not read response body: %w", res.StatusCode, err)
		}
		return fmt.Errorf("unexpected response '%d': %s", res.StatusCode, string(body))
	}

	return json.NewDecoder(bytes.NewReader(body)).Decode(output)
}
