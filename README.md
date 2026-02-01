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
*   **Model**: `openai/gpt-4o-mini`
*   **Authentication**: Requires a valid `GITHUB_TOKEN` in the environment variables.

### Key Capabilities

*   **Context Awareness**: The AI is prompted with specific contexts (e.g., Course Feedback, Employee Engagement) to tailor the analysis.
*   **Scalability**: The batching system ensures that large numbers of responses can be processed without hitting token limits.

## AI Usage Acknowledgment

This project utilizes Artificial Intelligence assistance, primarily for:

*   **Code Quality Enhancement**: Refactoring and optimizing code structure.
*   **Testing**: Generating comprehensive unit tests and increasing coverage.
*   **Documentation**: Drafting and refining documentation, including this README.
