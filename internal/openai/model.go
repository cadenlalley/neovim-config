package openai

import "gopkg.in/guregu/null.v4"

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string      `json:"role"`
			Content string      `json:"content"`
			Refusal null.String `json:"refusal"`
		} `json:"message"`
		LogProbs     null.String `json:"lobprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

type ChatCompletionContentImageURL struct {
	URL string `json:"url,omitempty"`
}

type ChatCompletionContent struct {
	Type     string                         `json:"type"`
	Text     string                         `json:"text,omitempty"`
	ImageURL *ChatCompletionContentImageURL `json:"image_url,omitempty"`
}

type ChatCompletionMessage struct {
	Role    string                  `json:"role"`
	Content []ChatCompletionContent `json:"content"`
}

type ChatCompletionRequest struct {
	Model          string                  `json:"model"`
	MaxTokens      int                     `json:"max_tokens"`
	Messages       []ChatCompletionMessage `json:"messages"`
	ResponseFormat interface{}             `json:"response_format"`
}
