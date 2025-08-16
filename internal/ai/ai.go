package ai

import (
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
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

type ResponseMetrics struct {
	Model            openai.ChatModel `json:"model"`
	PromptTokens     int64            `json:"promptTokens"`
	CompletionTokens int64            `json:"completionTokens"`
	Latency          int64            `json:"latency"`
}
