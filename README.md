# OSP Backend

## Table of Contents

- [Introduction](#introduction)
- [Tech Stack](#tech-stack)
- [How to run](#how-to-run)
  - [1) Configure environment variables](#1-configure-environment-variables)
  - [2) Start dependencies](#2-start-dependencies)
  - [3) Run the server](#3-run-the-server)
- [API Documentation](#api-documentation)
  - [Base URL](#base-url)
  - [Admin Authentication](#admin-authentication)
  - [Health Check](#health-check)
  - [Surveys](#surveys)
  - [Submissions](#submissions)
  - [Insights](#insights)
- [Example curl or HTTP requests for common flows](#example-curl-or-http-requests-for-common-flows)
- [AI Integration Details](#ai-integration-details)
  - [Architecture](#architecture)
  - [Providers & Models](#providers--models)
  - [Key Capabilities](#key-capabilities)
- [Future Work / Limitation](#future-work--limitation)
- [AI Usage Acknowledgment](#ai-usage-acknowledgment)

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
GET http://localhost:8080/api/health
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

#### Get Survey (Public)

Retrieve a survey by its public token.

- **GET** `/api/surveys/:token`
- Bruno: [.bruno/Get Survey.bru](.bruno/Get%20Survey.bru)

#### Create Survey (Admin)

Create a new survey with questions.

- **POST** `/api/admin/surveys`
- Bruno: [.bruno/Admin/Create Survey.bru](.bruno/Admin/Create%20Survey.bru)

Request Body Example:

```json
{
  "name": "Product Satisfaction Survey",
  "questions": [
    {
      "type": "TEXTBOX",
      "text": "What do you like most about our product?",
      "specification": {
        "max_length": 250
      }
    },
    {
      "type": "MULTIPLE_CHOICE",
      "text": "Which feature do you use most often?",
      "specification": {
        "options": [
          "Ease of use",
          "Price",
          "Customer support"
        ]
      }
    }
  ]
}
```

#### List Surveys (Admin)

List all surveys.

- **GET** `/api/admin/surveys`
- Bruno: [.bruno/Admin/List Surveys.bru](.bruno/Admin/List%20Surveys.bru)

#### Delete Survey (Admin)

- **DELETE** `/api/admin/surveys/:id`

---

### Submissions

#### Submit Survey (Public)

Submit responses for a survey.

- **POST** `/api/submissions`
- Bruno: [.bruno/Submit Survey.bru](.bruno/Submit%20Survey.bru)

Request Body Example:

```json
{
  "survey_token": "SURVEY_TOKEN",
  "responses": [
    {
      "question_id": "QUESTION_ID",
      "answer": "My detailed answer"
    }
  ]
}
```

#### List Submissions (Admin)

List submissions.

- **GET** `/api/admin/submissions`
- Bruno: [.bruno/Admin/Get Submissions.bru](.bruno/Admin/Get%20Submissions.bru)

#### Delete Submission (Admin)

- **DELETE** `/api/admin/submissions/:id`

---

### Insights

#### Create Insight (Admin)

Trigger AI insight generation for a survey.

- **POST** `/api/admin/insights`
- Bruno: [.bruno/Admin/Create Insight.bru](.bruno/Admin/Create%20Insight.bru)

Request Body Example:

```json
{
  "survey_id": "SURVEY_ID",
  "context_type": "PRODUCT_SATISFACTION"
}
```

#### List Insights (Admin)

- **GET** `/api/admin/insights`
- Bruno: [.bruno/Admin/Get Insights.bru](.bruno/Admin/Get%20Insights.bru)

#### Get Insight (Admin)

Retrieve a specific insight result.

- **GET** `/api/admin/insights/:id`
- Bruno: [.bruno/Admin/Get Insight.bru](.bruno/Admin/Get%20Insight.bru)

---

## Example curl or HTTP requests for common flows
Here is an example end-to-end flow using `curl`.

**1. Create a survey (Admin)**

```bash
curl -X POST http://localhost:8080/api/admin/surveys \
  -H "Authorization: Bearer root-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Product Satisfaction Survey",
    "questions": [
      {
        "type": "TEXTBOX",
        "text": "What did you like most?",
        "specification": { "max_length": 500 }
      }
    ]
  }'
```
*Response will contain the `id` (e.g., `67af9...`) and `token` (e.g., `550e8`).*

**2. Submit a response (Public)**

Use the `token` from the previous step.

```bash
curl -X POST http://localhost:8080/api/submissions \
  -H "Content-Type: application/json" \
  -d '{
    "survey_token": "550e8",
    "responses": [
      {
        "question_id": "qid_from_survey_creation",
        "answer": "The user interface is very clean."
      }
    ]
  }'
```

**3. Trigger Insight Generation (Admin)**

Use the `id` from step 1.

```bash
curl -X POST http://localhost:8080/api/admin/insights \
  -H "Authorization: Bearer root-token" \
  -H "Content-Type: application/json" \
  -d '{
    "survey_id": "67af9...",
    "context_type": "PRODUCT_SATISFACTION"
  }'
```
*Response will contain an `insight_id`.*

**4. Check Insight Result (Admin)**

Wait a few seconds for the background job to complete, then use the `insight_id`.

```bash
curl -X GET http://localhost:8080/api/admin/insights/insight_id_here \
  -H "Authorization: Bearer root-token"
```

## AI Integration Details

This project leverages Artificial Intelligence to generate actionable insights from survey responses.

### Architecture

1.  **Ingestion**: Survey submissions are collected via the API.
2.	**Batching**: Responses are grouped into batches to fit within context windows.
3.  **Queueing**: When an Insight is requested, a background job is enqueued using `asynq` (Redis).
4.  **Processing**:
    -   **Summarization**: Each batch is sent to the LLM for summarization.
    -   **Meta-Analysis**: All batch summaries are combined and sent for a final high-level analysis.
5.  **Storage**: Results are stored in MongoDB.

### Providers & Models

*   **Provider**: [GitHub Models](https://github.com/marketplace/models)
    *   **Free Tier Availability**: Utilizes a free tier suitable for initial prototyping usage without upfront costs.
    *   **Accessibility in Hong Kong**: Direct access supported in Hong Kong; no VPN is required, simplifying local development workflows.
    *   **GitHub Copilot Integration**: Leverages my existing GitHub Copilot subscription, which provides higher rate limits compared to standard free tier.
*   **Model**: `openai/gpt-4o-mini`
*   **Authentication**: Requires a valid `GITHUB_TOKEN` in the environment variables.

### Key Capabilities

*   **Scalability**: The batching system ensures that large numbers of responses can be processed without hitting token limits.
## Future Work / Limitation
*   **Test Verification**: Due to time constraints, currently only happy paths are tested, and not all AI-generated automated tests have been manually verified for edge cases.
*   **Survey Editing Strategy**: The ability to edit surveys is currently not implemented. Future support for this is planned to function as a "clone and modify" feature rather than in-place editing. This approach ensures data integrity by preserving the structure of surveys that already have submissions, avoiding inconsistencies between old responses and modified questions.
*   **Model Flexibility**: Support for multiple AI providers and models.
*   **Enhanced Prompting**: Dynamic prompt generation based on survey type.
*   **Error Handling**: The current error handling implementation is rudimentary and needs further improvement in the future

## AI Usage Acknowledgment

This project utilizes Artificial Intelligence assistance, primarily for:

*   **Code Quality Enhancement**: Refactoring and optimizing code structure.
*   **Testing**: Generating comprehensive unit tests and increasing coverage.
*   **Documentation**: Drafting and refining documentation, including this README.
