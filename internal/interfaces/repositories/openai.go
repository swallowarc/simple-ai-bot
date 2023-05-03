package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	openAIRepository struct {
		httpClient http.Client
		apiKey     string
	}

	ChatCompletionRequest struct {
		Model       string    `json:"model"`
		Messages    []Message `json:"messages"`
		MaxTokens   int       `json:"max_tokens"`
		Temperature float64   `json:"temperature"`
	}

	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	ChatCompletionResponse struct {
		ID      string    `json:"id"`
		Object  string    `json:"object"`
		Created int       `json:"created"`
		Model   string    `json:"model"`
		Usage   Usage     `json:"usage"`
		Choices []Choices `json:"choices"`
	}

	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}

	Choices struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	}
)

func NewOpenAIRepository(client http.Client) usecases.OpenAIRepository {
	return &openAIRepository{
		httpClient: client,
	}
}

func (c *openAIRepository) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
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

	client := &http.Client{}
	res, err := client.Do(request.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrapf(err, "error making request")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading response body")
	}

	var completionResponse ChatCompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return nil, errors.Wrapf(err, "error unmarshalling response body")
	}

	return &completionResponse, nil
}
