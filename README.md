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

## API Description

### Endpoints

#### Health Check
```
GET /api/health
```
Returns the health status of the API.

**Response:**
```json
{
  "status": "ok"
}
```

#### Surveys

##### Create New Survey
```
POST /api/surveys
```

**Request Body:**
```json
{
  "name": "test",
  "questions": [
    {
      "type": "TEXTBOX",
      "text": "Hi",
      "specification": {
        "max_length": 12
      }
    }
  ]
}
```

**Response:**
```json
{
  "id": "697e267a2dad33bddb80561b",
  "name": "test",
  "token": "IlXyz",
  "questions": [
    {
      "id": "697e267a2dad33bddb80561c",
      "text": "Hi",
      "type": "TEXTBOX",
      "specification": {
        "max_length": 12
      }
    }
  ],
  "created_at": "2026-01-31T23:57:46.2653345+08:00",
  "updated_at": "2026-01-31T23:57:46.2653345+08:00"
}
```