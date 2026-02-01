export interface QuestionSpecification {
  max_length?: number;
  options?: string[];
  min?: number;
  max?: number;
  min_label?: string;
  max_label?: string;
}

export interface Question {
  id?: string;
  type: 'TEXTBOX' | 'MULTIPLE_CHOICE' | 'LIKERT';
  text: string;
  specification: QuestionSpecification;
}

export interface Survey {
  id: string;
  token: string;
  name: string;
  questions: Question[];
  created_at: string;
  updated_at: string;
}

export interface CreateSurveyRequest {
  name: string;
  questions: Question[];
}

export interface SubmissionResponse {
  question_id: string;
  answer: string;
}

export interface Submission {
  id: string;
  survey_id: string;
  responses: SubmissionResponse[];
  created_at: string;
  updated_at: string;
}

export interface CreateSubmissionRequest {
  survey_token: string;
  responses: SubmissionResponse[];
}

export interface InsightBatch {
  batch_number: number;
  question: Question;
  summary: string;
}

export interface Insight {
  id: string;
  survey_id: string;
  context_type: string;
  status: 'PENDING' | 'PROCESSING' | 'COMPLETED' | 'FAILED';
  analysis: string;
  batches: InsightBatch[];
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface CreateInsightRequest {
  survey_id: string;
  context_type: string;
}
