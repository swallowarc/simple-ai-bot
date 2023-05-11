package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
)

type (
	ChatCompletionRequest struct {
		Model       string              `json:"model"`
		Messages    domain.ChatMessages `json:"messages"`
		MaxTokens   int                 `json:"max_tokens,omitempty"`
		Temperature float64             `json:"temperature,omitempty"`
	}

	ChatCompletionResponse struct {
		ID      string  `json:"id"`
		Object  string  `json:"object"`
		Created int     `json:"created"`
		Model   string  `json:"model"`
		Usage   Usage   `json:"usage"`
		Choices Choices `json:"choices"`
	}

	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}

	Choice struct {
		domain.ChatMessage `json:"message"`
		FinishReason       string `json:"finish_reason"`
		Index              int    `json:"index"`
	}
	Choices []Choice
)

func (cs Choices) Messages() domain.ChatMessages {
	l := len(cs)
	if l == 0 {
		return nil
	}

	var messages domain.ChatMessages
	for _, c := range cs {
		messages = append(messages, c.ChatMessage)
	}

	return messages
}

type (
	openAIClient struct {
		httpClient  *http.Client
		apiKey      string
		maxTokens   int
		temperature float64
	}

	Env struct {
		APIKey         string  `envconfig:"api_key" required:"true"`
		APIMaxTokens   int     `envconfig:"api_max_tokens" default:"400"`
		APITemperature float64 `envconfig:"api_temperature" default:"0.6"`
	}
)

func (c *openAIClient) ChatCompletion(ctx context.Context, messages domain.ChatMessages) (domain.ChatMessages, error) {
	if len(messages) == 0 {
		return nil, errors.New("chat messages is empty")
	}

	req := &ChatCompletionRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		MaxTokens:   c.maxTokens,
		Temperature: c.temperature,
	}
	jsonData, err := json.Marshal(*req)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling request data")
	}

	request, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "error creating request")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.httpClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrapf(err, "error making request")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading response body")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error status code: %d, body: %s", res.StatusCode, body)
	}

	var completionResponse ChatCompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return nil, errors.Wrapf(err, "error unmarshalling response body")
	}

	return completionResponse.Choices.Messages(), nil
}

func NewClient(httpClient *http.Client, env Env) interfaces.AIClient {
	return &openAIClient{
		httpClient:  httpClient,
		apiKey:      env.APIKey,
		maxTokens:   env.APIMaxTokens,
		temperature: env.APITemperature,
	}
}
