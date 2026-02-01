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
	QuestionTypeTextbox        QuestionType = "TEXTBOX"
	QuestionTypeMultipleChoice QuestionType = "MULTIPLE_CHOICE"
	QuestionTypeLikert         QuestionType = "LIKERT"
)

type QuestionSpecification struct {
	*TextboxSpecification        `bson:",inline,omitempty" json:",inline,omitempty"`
	*MultipleChoiceSpecification `bson:",inline,omitempty" json:",inline,omitempty"`
	*LikertSpecification         `bson:",inline,omitempty" json:",inline,omitempty"`
}

type TextboxSpecification struct {
	MaxLength int `bson:"max_length" json:"max_length" binding:"required,gt=0,lte=250"` // max 250 characters
}

type MultipleChoiceSpecification struct {
	Options []string `bson:"options" json:"options" binding:"required,min=2,max=20"`
}

type LikertSpecification struct {
	Min      int     `bson:"min" json:"min" binding:"required"`
	Max      int     `bson:"max" json:"max" binding:"required,gtfield=Min"`
	MinLabel *string `bson:"min_label" json:"min_label"`
	MaxLabel *string `bson:"max_label" json:"max_label"`
}

type QuestionInput struct {
	Type          QuestionType          `json:"type" binding:"required,oneof=TEXTBOX MULTIPLE_CHOICE LIKERT"`
	Text          string                `json:"text" binding:"required"`
	Specification QuestionSpecification `json:"specification" binding:"required"`
}

/* Request models */
type CreateSurveyRequest struct {
	Name      string          `json:"name" binding:"required"`
	Questions []QuestionInput `json:"questions" binding:"required,dive"`
}

type CreateSurveyResponse struct {
	Data  *Survey `json:"data"`
	Error string  `json:"error,omitempty"`
}

type ListSurveysRequest struct {
	Offset int64 `form:"offset,default=0"`
	Limit  int64 `form:"limit,default=10"`
}

type ListSurveysResponse struct {
	Data  []*Survey `json:"data"`
	Error string    `json:"error,omitempty"`
}

type GetSurveyByTokenRequest struct {
	Token string `uri:"token" binding:"required"`
}

type GetSurveyByTokenResponse struct {
	Data  *Survey `json:"data"`
	Error string  `json:"error,omitempty"`
}

type GetSurveyRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetSurveyResponse struct {
	Data  *Survey `json:"data"`
	Error string  `json:"error,omitempty"`
}

type DeleteSurveyRequest struct {
	ID string `uri:"id" binding:"required"`
}

type DeleteSurveyResponse struct {
	Error string `json:"error,omitempty"`
}
