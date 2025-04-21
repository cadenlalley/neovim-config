package ai

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AIClient struct {
	client openai.Client
}

func NewClient(token, host string) *AIClient {
	return &AIClient{
		client: openai.NewClient(
			option.WithAPIKey(token),
			option.WithBaseURL(host+"/v1"),
		),
	}
}

type TokenUsage struct {
	PromptTokens     int64 `json:"promptTokens"`
	CompletionTokens int64 `json:"completionTokens"`
	TotalTokens      int64 `json:"totalTokens"`
}
