# OSP Backend

## Table of Contents

- [Introduction](#introduction)
- [Tech Stack](#tech-stack)
- [API Documentation](#api-documentation)

## Introduction

## Tech Stack

- **Golang** - Backend programming language
- **Gin** - HTTP web framework
- **MongoDB** - NoSQL database
- **GitHub Models API** - AI integration

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
			"type": "TEXTBOX",
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
