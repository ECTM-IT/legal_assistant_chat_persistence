# Subscriptions API Documentation

## Overview
The Subscriptions API provides endpoints for managing user subscriptions, including subscription creation, updates, and payment processing through Stripe integration. It handles subscription lifecycles, billing periods, and payment status.

## Base URL
All endpoints are relative to the base URL of your API server.

## Authentication
All endpoints require authentication. Include your authentication token in the Authorization header:
```
Authorization: Bearer <your_token>
```

## Endpoints

### Get All Subscriptions
Retrieves all subscriptions.

**Endpoint:** `GET /subscriptions/`

**Response:** `200 OK`
```json
[
  {
    "id": "string",
    "user_id": "string",
    "plan": {
      "name": "string",
      "type": "string",
      "price": 0.0,
      "description": "string",
      "features": ["string"]
    },
    "expiry": "2024-01-01T00:00:00Z",
    "status": "string",
    "stripe_customer_id": "string",
    "stripe_subscription_id": "string",
    "current_period_start": "2024-01-01T00:00:00Z",
    "current_period_end": "2024-01-01T00:00:00Z",
    "cancel_at_period_end": boolean,
    "billing_informations": {}
  }
]
```

### Get Subscription by ID
Retrieves a specific subscription by its ID.

**Endpoint:** `GET /subscriptions/{id}/`

**Path Parameters:**
- `id`: Subscription ID (string)

**Response:** `200 OK`
```json
{
  "id": "string",
  "user_id": "string",
  "plan": {
    "name": "string",
    "type": "string",
    "price": 0.0,
    "description": "string",
    "features": ["string"]
  },
  "expiry": "2024-01-01T00:00:00Z",
  "status": "string",
  "stripe_customer_id": "string",
  "stripe_subscription_id": "string",
  "current_period_start": "2024-01-01T00:00:00Z",
  "current_period_end": "2024-01-01T00:00:00Z",
  "cancel_at_period_end": boolean,
  "billing_informations": {}
}
```

### Get Subscriptions by Plan
Retrieves all subscriptions for a specific plan.

**Endpoint:** `GET /subscriptions/plan/`

**Query Parameters:**
- `plan_name`: Plan name (string)
- `plan_type`: Plan type (string)

**Response:** `200 OK`
```json
[
  {
    "id": "string",
    "user_id": "string",
    "plan": {
      "name": "string",
      "type": "string",
      "price": 0.0,
      "description": "string",
      "features": ["string"]
    },
    "expiry": "2024-01-01T00:00:00Z",
    "status": "string",
    "stripe_customer_id": "string",
    "stripe_subscription_id": "string",
    "current_period_start": "2024-01-01T00:00:00Z",
    "current_period_end": "2024-01-01T00:00:00Z",
    "cancel_at_period_end": boolean,
    "billing_informations": {}
  }
]
```

### Create Subscription
Creates a new subscription.

**Endpoint:** `POST /subscriptions/`

**Request Body:**
```json
{
  "user_id": "string",
  "plan": {
    "name": "string",
    "type": "string"
  }
}
```

**Response:** `201 Created`
```json
{
  "id": "string",
  "user_id": "string",
  "plan": {
    "name": "string",
    "type": "string",
    "price": 0.0,
    "description": "string",
    "features": ["string"]
  },
  "expiry": "2024-01-01T00:00:00Z",
  "status": "string",
  "stripe_customer_id": "string",
  "stripe_subscription_id": "string",
  "current_period_start": "2024-01-01T00:00:00Z",
  "current_period_end": "2024-01-01T00:00:00Z",
  "cancel_at_period_end": false,
  "billing_informations": {}
}
```

### Update Subscription
Updates an existing subscription.

**Endpoint:** `PATCH /subscriptions/{id}/`

**Path Parameters:**
- `id`: Subscription ID (string)

**Request Body:**
```json
{
  "plan": {
    "name": "string",
    "type": "string"
  },
  "status": "string",
  "cancel_at_period_end": boolean
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "user_id": "string",
  "plan": {
    "name": "string",
    "type": "string",
    "price": 0.0,
    "description": "string",
    "features": ["string"]
  },
  "expiry": "2024-01-01T00:00:00Z",
  "status": "string",
  "stripe_customer_id": "string",
  "stripe_subscription_id": "string",
  "current_period_start": "2024-01-01T00:00:00Z",
  "current_period_end": "2024-01-01T00:00:00Z",
  "cancel_at_period_end": boolean,
  "billing_informations": {}
}
```

### Delete Subscription
Cancels and deletes a subscription.

**Endpoint:** `DELETE /subscriptions/{id}/`

**Path Parameters:**
- `id`: Subscription ID (string)

**Response:** `204 No Content`

### Purchase Subscription
Processes a subscription purchase.

**Endpoint:** `POST /subscriptions/purchase/`

**Request Body:**
```json
{
  "plan_name": "string",
  "plan_type": "string",
  "payment_method_id": "string",
  "billing_informations": {}
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "user_id": "string",
  "plan": {
    "name": "string",
    "type": "string",
    "price": 0.0,
    "description": "string",
    "features": ["string"]
  },
  "expiry": "2024-01-01T00:00:00Z",
  "status": "string",
  "stripe_customer_id": "string",
  "stripe_subscription_id": "string",
  "current_period_start": "2024-01-01T00:00:00Z",
  "current_period_end": "2024-01-01T00:00:00Z",
  "cancel_at_period_end": false,
  "billing_informations": {}
}
```

## Data Types

### Subscription Object
```json
{
  "id": "string",
  "user_id": "string",
  "plan": {
    "name": "string",
    "type": "string",
    "price": 0.0,
    "description": "string",
    "features": ["string"]
  },
  "expiry": "2024-01-01T00:00:00Z",
  "status": "string",
  "stripe_customer_id": "string",
  "stripe_subscription_id": "string",
  "current_period_start": "2024-01-01T00:00:00Z",
  "current_period_end": "2024-01-01T00:00:00Z",
  "cancel_at_period_end": boolean,
  "billing_informations": {}
}
```

### Subscription Status
The status field can have the following values:
- `active`
- `canceled`
- `past_due`

### Validation Rules
- User ID is required
- Plan name and type must be valid according to the Plans API
- Stripe customer ID is required for active subscriptions
- Stripe subscription ID is required for active subscriptions
- Current period start and end dates must be valid timestamps
- Status must be one of the defined status values
- Billing information must be a valid JSON object

## Error Responses

### Common HTTP Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `204 No Content`: Request successful, no content to return
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Subscription not found
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