package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

/* Main models */
type ChatCompletionRequestLog struct {
	ID        bson.ObjectID           `bson:"_id" json:"id"`
	Request   ChatCompletionRequest   `bson:"request" json:"request"`
	Response  *ChatCompletionResponse `bson:"response,omitempty" json:"response,omitempty"`
	Reference *string                 `bson:"reference" json:"reference"`
	CreatedAt time.Time               `bson:"created_at" json:"created_at"`
}

// LLM request/response structures type Message struct { Role string `json:"role"` Content string `json:"content"` }
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Messages    []ChatCompletionMessage `json:"messages"`
	Temperature float64                 `json:"temperature"`
	TopP        float64                 `json:"top_p"`
	MaxTokens   int                     `json:"max_tokens"`
	Model       string                  `json:"model"`
}

type ChatCompletionChoice struct {
	Message ChatCompletionMessage `json:"message"`
}

type ChatCompletionResponse struct {
	Choices []ChatCompletionChoice `json:"choices"`
}
