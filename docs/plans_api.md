# Plans API Documentation

## Overview
The Plans API provides endpoints for managing subscription plans, including retrieving plan options, toggling between monthly and annual plans, and selecting plans. The system offers three tiers: Pro, Team, and Enterprise, each available in monthly and annual billing cycles.

## Base URL
All endpoints are relative to the base URL of your API server.

## Authentication
All endpoints require authentication. Include your authentication token in the Authorization header:
```
Authorization: Bearer <your_token>
```

## Endpoints

### Get Plan Options
Retrieves all available plan options.

**Endpoint:** `GET /plans/`

**Response:** `200 OK`
```json
{
  "monthly": [
    {
      "name": "pro",
      "type": "monthly",
      "price": 20.00,
      "description": "Professional plan for individual users",
      "features": [
        "Unlimited cases",
        "Priority support",
        "Advanced analytics"
      ]
    },
    {
      "name": "team",
      "type": "monthly",
      "price": 60.00,
      "description": "Perfect for small teams",
      "features": [
        "Everything in Pro",
        "Team collaboration",
        "Admin dashboard"
      ]
    },
    {
      "name": "enterprise",
      "type": "monthly",
      "price": 0,
      "description": "For large organizations",
      "features": [
        "Everything in Team",
        "Custom integrations",
        "Dedicated support"
      ]
    }
  ],
  "annual": [
    {
      "name": "pro",
      "type": "annual",
      "price": 192.00,
      "description": "Professional plan for individual users",
      "features": [
        "Unlimited cases",
        "Priority support",
        "Advanced analytics"
      ]
    },
    {
      "name": "team",
      "type": "annual",
      "price": 576.00,
      "description": "Perfect for small teams",
      "features": [
        "Everything in Pro",
        "Team collaboration",
        "Admin dashboard"
      ]
    },
    {
      "name": "enterprise",
      "type": "annual",
      "price": 0,
      "description": "For large organizations",
      "features": [
        "Everything in Team",
        "Custom integrations",
        "Dedicated support"
      ]
    }
  ]
}
```

### Toggle Plan Type
Toggles between monthly and annual billing cycles.

**Endpoint:** `PATCH /plans/toggle/`

**Response:** `200 OK`
```json
{
  "type": "string" // "monthly" or "annual"
}
```

### Select Plan
Selects a specific plan for the user.

**Endpoint:** `POST /plans/select/`

**Request Body:**
```json
{
  "name": "string",
  "type": "string"
}
```

**Response:** `200 OK`
```json
{
  "name": "string",
  "type": "string",
  "price": 0.0,
  "description": "string",
  "features": ["string"]
}
```

## Data Types

### Plan Object
```json
{
  "name": "string",
  "type": "string",
  "price": 0.0,
  "description": "string",
  "features": ["string"]
}
```

### Plan Types
The type field can have the following values:
- `monthly`
- `annual`

### Plan Names
The name field can have the following values:
- `pro`
- `team`
- `enterprise`

### Validation Rules
- Plan name must be one of the predefined values
- Plan type must be either "monthly" or "annual"
- Price must be a non-negative number
- Description is required
- Features must be a non-empty array of strings

## Error Responses

### Common HTTP Status Codes
- `200 OK`: Request successful
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Plan not found
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