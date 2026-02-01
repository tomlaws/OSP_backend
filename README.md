# OSP Backend

## Table of Contents

- [Introduction](#introduction)
- [Tech Stack](#tech-stack)
- [How to run](#how-to-run)
- [API Documentation](#api-documentation)

## Introduction
OSP Backend is the backend service for the OSP (Open Survey Platform) application. It provides RESTful APIs for managing surveys, collecting submissions, and generating insights using AI.

## Tech Stack

- **Golang** - Backend programming language
- **Gin** - HTTP web framework
- **MongoDB** - NoSQL database
- **Redis + Asynq** - Background job processing
- **GitHub Models API** - AI integration

## How to run

### 1) Configure environment variables

Create a `.env` file in the project root (same folder as `go.mod`) and set:

```env
PORT=8080

# Admin Bearer token for /api/admin/* endpoints
ROOT_TOKEN=change-me

# MongoDB
MONGODB_URI=mongodb://localhost:27017

# Redis (required by Asynq job system)
REDIS_URI=localhost:6379

# Required for insights generation via GitHub Models API
GITHUB_TOKEN=your_github_token
```

Notes:

- If `ROOT_TOKEN` is empty, admin endpoints will return `401 Unauthorized`.
- `GITHUB_TOKEN` is only needed for accessing the GitHub Models API for insights generation.

### 2) Start dependencies

- Start MongoDB (local or hosted)
- Start Redis (for Asynq)

### 3) Run the server

From the repo root:

```bash
go run ./cmd/server
```

Then verify:

```text
GET http://localhost:8080/health
```

## API Documentation

The quickest way to try the API is via the Bruno collection in [.bruno](.bruno). The requests shown below are based on those `.bru` files.

### Base URL

Bruno environments define `BASE_URL` (see [.bruno/environments/Local.bru](.bruno/environments/Local.bru)).

Example:

```text
http://localhost:8080
```

### Admin Authentication

All `/api/admin/*` endpoints require:

```text
Authorization: Bearer <ROOT_TOKEN>
```

This is the same token configured by `ROOT_TOKEN` in your environment.

---

### Health Check

- **GET** `/health`
- Bruno: [.bruno/Health Check.bru](.bruno/Health%20Check.bru)

Sample response:

```json
{
	"status": "ok"
}
```

---

### Surveys

#### Create Survey

- **POST** `/api/surveys`
- Bruno: [.bruno/Create Survey.bru](.bruno/Create%20Survey.bru)

Sample request:

```json
{
	"name": "Product Satisfaction Survey",
	"questions": [
		{
			"type": "TEXTBOX", // Possible values: TEXTBOX, MULTIPLE_CHOICE, LIKERT
			"text": "What do you like most about our product?",
			"specification": { "max_length": 250 }
		},
		{
			"type": "MULTIPLE_CHOICE",
			"text": "Which feature do you use most often?",
			"specification": {
				"options": ["Ease of use", "Price", "Customer support", "Design", "Performance"]
			}
		},
		{
			"type": "LIKERT",
			"text": "How satisfied are you with the product overall?",
			"specification": {
				"min": 1,
				"max": 5,
				"minLabel": "Very dissatisfied",
				"maxLabel": "Very satisfied"
			}
		}
	]
}
```

Sample response (shape):

```json
{
	"id": "697ec2067cd24f1b1553146e",
	"name": "Product Satisfaction Survey",
	"token": "YUNvS",
	"questions": [
		{
			"id": "697ec2067cd24f1b1553146f",
			"text": "What do you like most about our product?",
			"type": "TEXTBOX",
			"specification": { "max_length": 250 }
		}
	],
	"created_at": "2026-02-01T12:00:00Z",
	"updated_at": "2026-02-01T12:00:00Z"
}
```

#### Get Survey

- **GET** `/api/surveys/:token`
- Bruno: [.bruno/Get Survey.bru](.bruno/Get%20Survey.bru)

Sample response: same shape as **Create Survey**.

---

### Submissions

#### Submit Survey

- **POST** `/api/submissions`
- Bruno: [.bruno/Submit Survey.bru](.bruno/Submit%20Survey.bru)

Sample request:

```json
{
	"survey_token": "YUNvS",
	"responses": [
		{ "question_id": "697ec2067cd24f1b1553146f", "answer": "Attractive design" },
		{ "question_id": "697ec2067cd24f1b15531470", "answer": "Design" },
		{ "question_id": "697ec2067cd24f1b15531471", "answer": "5" }
	]
}
```

Sample response:

```json
{
	"data": {
		"id": "697ec2067cd24f1b15539999",
		"survey_id": "697ec2067cd24f1b1553146e",
		"responses": [
			{ "question_id": "697ec2067cd24f1b1553146f", "answer": "Attractive design" }
		],
		"created_at": "2026-02-01T12:05:00Z",
		"updated_at": "2026-02-01T12:05:00Z"
	}
}
```

---

### Insights (Admin)

#### Create Insight

- **POST** `/api/admin/insights`
- Auth: Bearer `ROOT_TOKEN`
- Bruno: [.bruno/Create Insight.bru](.bruno/Create%20Insight.bru)

Sample request:

```json
{
	"survey_id": "697ec2067cd24f1b1553146e",
	"context_type": "PRODUCT_SATISFACTION"
}
```

Sample response (shape):

```json
{
	"id": "697ec2067cd24f1b15540000",
	"survey_id": "697ec2067cd24f1b1553146e",
	"context_type": "PRODUCT_SATISFACTION",
	"status": "PENDING",
	"analysis": "",
	"batches": [
		{
			"batch_number": 1,
			"question": { "id": "697ec2067cd24f1b1553146f", "text": "...", "type": "TEXTBOX", "specification": { "max_length": 250 } },
			"textual_answers": [],
			"summary": null,
			"error_log": null
		}
	],
	"created_at": "2026-02-01T12:10:00Z",
	"updated_at": "2026-02-01T12:10:00Z",
	"completed_at": null
}
```

#### Get Insights

- **GET** `/api/admin/insights?offset=0&limit=10&surveyId=<surveyId>`
- Auth: Bearer `ROOT_TOKEN`
- Bruno: [.bruno/Get Insights.bru](.bruno/Get%20Insights.bru)
- Filtered example: [.bruno/Get Insight.bru](.bruno/Get%20Insight.bru)

Sample response:

```json
[
	{
		"id": "697ec2067cd24f1b15540000",
		"survey_id": "697ec2067cd24f1b1553146e",
		"context_type": "PRODUCT_SATISFACTION",
		"status": "COMPLETED",
		"analysis": "Overall analysis...",
		"batches": [],
		"created_at": "2026-02-01T12:10:00Z",
		"updated_at": "2026-02-01T12:12:00Z",
		"completed_at": "2026-02-01T12:12:00Z"
	}
]
```

Notes:

- Results are sorted by `completed_at` descending.
