package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"osp/internal/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

const githubModelsChatCompletionsURL = "https://models.github.ai/inference/chat/completions"

type ChatCompletionService struct {
	collection *mongo.Collection
}

func NewChatCompletionService(collection *mongo.Collection) *ChatCompletionService {
	return &ChatCompletionService{
		collection: collection,
	}
}

func (s *ChatCompletionService) NewRequest(reqBody models.ChatCompletionRequest) (*string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is not set")
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, githubModelsChatCompletionsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("github models request failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var chatCompletionResponse models.ChatCompletionResponse
	if err := json.Unmarshal(body, &chatCompletionResponse); err != nil {
		return nil, err
	}

	c, err := firstChoiceContent(&chatCompletionResponse)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func firstChoiceContent(resp *models.ChatCompletionResponse) (string, error) {
	if resp == nil || len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}
	return resp.Choices[0].Message.Content, nil
}
