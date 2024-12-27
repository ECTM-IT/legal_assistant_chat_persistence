# Users API Documentation

## Overview
The Users API provides endpoints for managing user accounts, including user creation, updates, and retrieval. Users can be associated with cases, teams, agents, and subscriptions.

## Base URL
All endpoints are relative to the base URL of your API server.

## Authentication
All endpoints require authentication. Include your authentication token in the Authorization header:
```
Authorization: Bearer <your_token>
```

## Endpoints

### Get User by ID
Retrieves a specific user by their ID.

**Endpoint:** `GET /users/{id}/`

**Path Parameters:**
- `id`: User ID (string)

**Response:** `200 OK`
```json
{
  "id": "string",
  "image": "string",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string",
  "case_ids": ["string"],
  "team_id": "string",
  "agent_ids": ["string"],
  "subscription_id": "string"
}
```

### Get User by Email
Retrieves a user by their email address.

**Endpoint:** `POST /users-email/`

**Request Body:**
```json
{
  "email": "string"
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "image": "string",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string",
  "case_ids": ["string"],
  "team_id": "string",
  "agent_ids": ["string"],
  "subscription_id": "string"
}
```

### Create User
Creates a new user account.

**Endpoint:** `POST /users/`

**Request Body:**
```json
{
  "image": "string",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string"
}
```

**Response:** `201 Created`
```json
{
  "id": "string",
  "image": "string",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string",
  "case_ids": [],
  "team_id": "string",
  "agent_ids": [],
  "subscription_id": "string"
}
```

### Update User
Updates an existing user's information.

**Endpoint:** `PATCH /users/{id}/`

**Path Parameters:**
- `id`: User ID (string)

**Request Body:**
```json
{
  "image": "string",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string",
  "case_ids": ["string"],
  "team_id": "string",
  "agent_ids": ["string"],
  "subscription_id": "string"
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "image": "string",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string",
  "case_ids": ["string"],
  "team_id": "string",
  "agent_ids": ["string"],
  "subscription_id": "string"
}
```

### Delete User
Deletes a user account.

**Endpoint:** `DELETE /users/{id}/`

**Path Parameters:**
- `id`: User ID (string)

**Response:** `204 No Content`

## Data Types

### User Object
```json
{
  "id": "string",
  "image": "string",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string",
  "case_ids": ["string"],
  "team_id": "string",
  "agent_ids": ["string"],
  "subscription_id": "string"
}
```

### Validation Rules
- Email is required and must be a valid email format
- First name is required
- Last name is required
- Phone number must be in a valid format (if provided)
- Image URL must be a valid URL (if provided)
- Case IDs must be valid ObjectIDs
- Team ID must be a valid ObjectID
- Agent IDs must be valid ObjectIDs
- Subscription ID must be a valid ObjectID

## Error Responses

### Common HTTP Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `204 No Content`: Request successful, no content to return
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: User not found
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