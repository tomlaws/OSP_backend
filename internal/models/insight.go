package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

/* Main models */
type Insight struct {
	ID             bson.ObjectID `bson:"_id" json:"id"`
	SurveyID       bson.ObjectID `bson:"survey_id" json:"survey_id"`
	ContextType    ContextType   `bson:"context_type" json:"context_type"`
	Status         InsightStatus `bson:"status" json:"status"`
	Analysis       string        `bson:"analysis" json:"analysis"`
	Batches        []AnswerBatch `bson:"batches"`
	BatchSummaries []string      `bson:"batch_summaries"`
	CurrentBatch   int           `bson:"current_batch"`
	ErrorLog       []string      `bson:"error_log,omitempty"` // track failed batches
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
	CompletedAt    *time.Time    `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

type AnswerBatch struct {
	BatchNumber      int            `bson:"batch_number"`
	Question         string         `bson:"question" json:"question"`
	AggregatedAnswer map[string]int `bson:"aggregated_answer,omitempty"`
	TextualAnswers   []string       `bson:"textual_answers,omitempty"`
	Summary          string         `bson:"summary,omitempty"` // store LLM output directly
}

type ContextType string

const (
	CourseFeedbackContext      ContextType = "COURSE_FEEDBACK"
	ProductSatisfactionContext ContextType = "PRODUCT_SATISFACTION"
	EmployeeEngagementContext  ContextType = "EMPLOYEE_ENGAGEMENT"
	EventFeedbackContext       ContextType = "EVENT_FEEDBACK"
)

type InsightStatus string

const (
	InsightPending    InsightStatus = "PENDING"
	InsightProcessing InsightStatus = "PROCESSING"
	InsightCompleted  InsightStatus = "COMPLETED"
	InsightFailed     InsightStatus = "FAILED"
)

/* Request models */
type CreateInsightRequest struct {
	SurveyID    bson.ObjectID `json:"survey_id" binding:"required"`
	ContextType ContextType   `json:"context_type" binding:"required,oneof=COURSE_FEEDBACK PRODUCT_SATISFACTION EMPLOYEE_ENGAGEMENT EVENT_FEEDBACK"`
}
