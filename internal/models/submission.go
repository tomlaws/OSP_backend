package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

/* Main models */
type Submission struct {
	ID        bson.ObjectID        `bson:"_id" json:"id"`
	SurveyID  bson.ObjectID        `bson:"survey_id" json:"survey_id" binding:"required"`
	Responses []SubmissionResponse `bson:"responses" json:"responses" binding:"required"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

type SubmissionResponse struct {
	QuestionID bson.ObjectID `bson:"question_id" json:"question_id" binding:"required"`
	Answer     string        `bson:"answer" json:"answer" binding:"required"`
}

/* Request models */
type CreateSubmissionRequest struct {
	SurveyToken string               `json:"survey_token" binding:"required"`
	Responses   []SubmissionResponse `json:"responses" binding:"required"`
}

/* Response models */
type CreateSubmissionResponse struct {
	Data *Submission `json:"data"`
}
