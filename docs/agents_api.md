# Agents API Documentation

## Overview
The Agents API provides endpoints for managing AI agents in the system. It allows you to retrieve agent information, get agents by user, and handle agent purchases.

## Base URL
All endpoints are relative to the base URL of your API server.

## Authentication
All endpoints require authentication. Include your authentication token in the Authorization header:
```
Authorization: Bearer <your_token>
```

## Endpoints

### Get All Agents
Retrieves all available agents in the system.

**Endpoint:** `GET /agents/`

**Response:** `200 OK`
```json
[
  {
    "id": "string",
    "profile_image": "string",
    "name": "string",
    "description": "string",
    "skills": ["string"],
    "price": 0.0,
    "code": "string"
  }
]
```

### Get Agent by ID
Retrieves a specific agent by its ID.

**Endpoint:** `GET /agents/{id}/`

**Path Parameters:**
- `id`: Agent ID (string)

**Response:** `200 OK`
```json
{
  "id": "string",
  "profile_image": "string",
  "name": "string",
  "description": "string",
  "skills": ["string"],
  "price": 0.0,
  "code": "string"
}
```

### Get Agents by User
Retrieves all agents associated with the authenticated user.

**Endpoint:** `GET /agents-user/`

**Response:** `200 OK`
```json
[
  {
    "id": "string",
    "profile_image": "string",
    "name": "string",
    "description": "string",
    "skills": ["string"],
    "price": 0.0,
    "code": "string"
  }
]
```

### Purchase Agent
Purchases an agent for the authenticated user.

**Endpoint:** `GET /agent-purchase/{id}/`

**Path Parameters:**
- `id`: Agent ID (string)

**Response:** `200 OK`
```json
{
  "id": "string",
  "profile_image": "string",
  "name": "string",
  "description": "string",
  "skills": ["string"],
  "price": 0.0,
  "code": "string"
}
```

## Data Types

### Agent Object
```json
{
  "id": "string",
  "profile_image": "string",
  "name": "string",
  "description": "string",
  "skills": ["string"],
  "price": 0.0,
  "code": "string"
}
```

### Validation Rules
- Agent name is required
- Price must be a positive number
- Code is required and must be unique
- Skills must be an array of strings

## Error Responses

### Common HTTP Status Codes
- `200 OK`: Request successful
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Agent not found
- `422 Unprocessable Entity`: Validation error
- `500 Internal Server Error`: Server error

### Error Response Format
```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": {}
  }
}
``` 