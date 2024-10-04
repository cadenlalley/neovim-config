package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

func NewOpenAIClient(baseURL, apiKey string) *OpenAIClient {
	return &OpenAIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *OpenAIClient) PostChatCompletion(payload ChatCompletionRequest) (ChatCompletionResponse, error) {
	// Construct the URL
	reqURL := c.BaseURL + "/v1/chat/completions"

	// Marshal the payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(data))
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	defer resp.Body.Close()

	// Handle errors
	if resp.StatusCode != http.StatusOK {
		return ChatCompletionResponse{}, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	var result ChatCompletionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	return result, nil
}
