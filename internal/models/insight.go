package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

/* Main models */
type Insight struct {
	ID          bson.ObjectID  `bson:"_id" json:"id"`
	SurveyID    bson.ObjectID  `bson:"survey_id" json:"survey_id"`
	ContextType ContextType    `bson:"context_type" json:"context_type"`
	Status      InsightStatus  `bson:"status" json:"status"`
	Analysis    string         `bson:"analysis" json:"analysis"`
	Batches     []InsightBatch `bson:"batches" json:"batches"`
	CreatedAt   time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at" json:"updated_at"`
	CompletedAt *time.Time     `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

type InsightBatch struct {
	BatchNumber      int             `bson:"batch_number" json:"batch_number"`
	Question         Question        `bson:"question" json:"question"`
	AggregatedAnswer *map[string]int `bson:"aggregated_answer,omitempty" json:"aggregated_answer,omitempty"`
	TextualAnswers   *[]string       `bson:"textual_answers,omitempty" json:"textual_answers,omitempty"`
	Summary          *string         `bson:"summary,omitempty" json:"summary,omitempty"`
	ErrorLog         *string         `bson:"error_log,omitempty" json:"error_log,omitempty"`
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

type CreateInsightResponse struct {
	Data  *Insight `json:"data"`
	Error string   `json:"error,omitempty"`
}

type GetInsightsRequest struct {
	Offset   int64   `form:"offset,default=0"`
	Limit    int64   `form:"limit,default=10"`
	SurveyID *string `form:"surveyId"`
}

type GetInsightsResponse struct {
	Data  []*Insight `json:"data"`
	Error string     `json:"error,omitempty"`
}

type GetInsightRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetInsightResponse struct {
	Data  *Insight `json:"data"`
	Error string   `json:"error,omitempty"`
}
