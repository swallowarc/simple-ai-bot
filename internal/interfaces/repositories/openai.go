package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"

	"github.com/pkg/errors"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	openAIRepository struct {
		httpClient *http.Client
		apiKey     string
	}

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

func NewOpenAIRepository(client *http.Client, apiKey string) usecases.OpenAIRepository {
	return &openAIRepository{
		httpClient: client,
		apiKey:     apiKey,
	}
}

func (r *openAIRepository) ChatCompletion(ctx context.Context, messages domain.ChatMessages) (domain.ChatMessages, error) {
	req := &ChatCompletionRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
		//Temperature: 0,
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
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.apiKey))

	res, err := r.httpClient.Do(request.WithContext(ctx))
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
