package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

/* Main models */
type Survey struct {
	ID        bson.ObjectID `bson:"_id" json:"id" binding:"required"`
	Name      string        `bson:"name" json:"name" binding:"required"`
	Token     string        `bson:"token" json:"token" binding:"required"`
	Questions []Question    `bson:"questions" json:"questions" binding:"required"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type Question struct {
	ID            bson.ObjectID         `bson:"id" json:"id" binding:"required"`
	Text          string                `bson:"text" json:"text" binding:"required"`
	Type          QuestionType          `bson:"type" json:"type" binding:"required"`
	Specification QuestionSpecification `bson:"specification" json:"specification"`
}

type QuestionType string

const (
	Textbox        QuestionType = "TEXTBOX"
	MultipleChoice QuestionType = "MULTIPLE_CHOICE"
	Likert         QuestionType = "LIKERT"
)

type QuestionSpecification struct {
	*TextboxSpecification        `bson:",inline,omitempty" json:",inline,omitempty"`
	*MultipleChoiceSpecification `bson:",inline,omitempty" json:",inline,omitempty"`
	*LikertSpecification         `bson:",inline,omitempty" json:",inline,omitempty"`
}

type TextboxSpecification struct {
	MaxLength int `bson:"max_length" json:"max_length"`
}

type MultipleChoiceSpecification struct {
	Options []string `bson:"options" json:"options"`
}

type LikertSpecification struct {
	Min      int    `bson:"min" json:"min"`
	Max      int    `bson:"max" json:"max"`
	MinLabel string `bson:"min_label" json:"min_label"`
	MaxLabel string `bson:"max_label" json:"max_label"`
}

/* Request models */
type CreateSurveyRequest struct {
	Name      string          `json:"name" binding:"required"`
	Questions []QuestionInput `json:"questions" binding:"required,dive"`
}

type QuestionInput struct {
	Type          QuestionType          `json:"type" binding:"required,oneof=TEXTBOX MULTIPLE_CHOICE LIKERT"`
	Text          string                `json:"text" binding:"required"`
	Specification QuestionSpecification `json:"specification" binding:"required"`
}

/* Response models */
type GetSurveyResponse struct {
	Data *Survey `json:"data"`
}
