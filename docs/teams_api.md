# Teams API Documentation

## Overview
The Teams API provides endpoints for managing teams, team members, and team invitations. It allows you to create and manage teams, add/update/remove team members, and handle team invitations. The API implements soft deletion for both teams and team members.

## Base URL
All endpoints are relative to the base URL of your API server.

## Authentication
All endpoints require authentication. Include your authentication token in the Authorization header:
```
Authorization: Bearer <your_token>
```

## Endpoints

### Create Team
Creates a new team with optional members.

**Endpoint:** `POST /teams/`

**Request Body:**
```json
{
  "name": "string",
  "description": "string",
  "members": [
    {
      "email": "string",
      "user_id": "string",
      "role": "string",
      "first_name": "string",
      "last_name": "string"
    }
  ]
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "members": [
    {
      "id": "string",
      "user_id": "string",
      "role": "string",
      "first_name": "string",
      "last_name": "string",
      "email": "string",
      "date_added": "2024-01-01T00:00:00Z",
      "last_active": "2024-01-01T00:00:00Z",
      "is_deleted": false,
      "deleted_at": null
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "is_deleted": false,
  "deleted_at": null
}
```

### Get Team by ID
Retrieves a specific team by its ID.

**Endpoint:** `GET /teams/{id}/`

**Path Parameters:**
- `id`: Team ID (string)

**Response:** `200 OK`
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "members": [
    {
      "id": "string",
      "user_id": "string",
      "role": "string",
      "first_name": "string",
      "last_name": "string",
      "email": "string",
      "date_added": "2024-01-01T00:00:00Z",
      "last_active": "2024-01-01T00:00:00Z",
      "is_deleted": false,
      "deleted_at": null
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "is_deleted": false,
  "deleted_at": null
}
```

### Get All Teams
Retrieves all teams accessible to the authenticated user.

**Endpoint:** `GET /teams/`

**Response:** `200 OK`
```json
[
  {
    "id": "string",
    "name": "string",
    "description": "string",
    "members": [...],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### Update Team
Updates an existing team's information.

**Endpoint:** `PUT /teams/{id}/`

**Path Parameters:**
- `id`: Team ID (string)

**Request Body:**
```json
{
  "name": "string",
  "description": "string"
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "members": [...],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "is_deleted": false,
  "deleted_at": null
}
```

### Delete Team
Soft deletes a team. The team will be marked as deleted but not permanently removed from the database.

**Endpoint:** `DELETE /teams/{id}/`

**Path Parameters:**
- `id`: Team ID (string)

**Response:** `204 No Content`

### Restore Team
Restores a previously deleted team.

**Endpoint:** `POST /teams/{id}/restore/`

**Path Parameters:**
- `id`: Team ID (string)

**Response:** `200 OK`

### Add Team Member
Adds a new member to a team.

**Endpoint:** `POST /teams/{id}/members/`

**Path Parameters:**
- `id`: Team ID (string)

**Request Body:**
```json
{
  "email": "string",
  "role": "string",
  "first_name": "string",
  "last_name": "string"
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "user_id": "string",
  "role": "string",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "date_added": "2024-01-01T00:00:00Z",
  "last_active": "2024-01-01T00:00:00Z",
  "is_deleted": false,
  "deleted_at": null
}
```

### Update Team Member
Updates an existing team member's information.

**Endpoint:** `PUT /teams/{id}/members/{memberId}/`

**Path Parameters:**
- `id`: Team ID (string)
- `memberId`: Member ID (string)

**Request Body:**
```json
{
  "role": "string",
  "first_name": "string",
  "last_name": "string",
  "email": "string"
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "user_id": "string",
  "role": "string",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "date_added": "2024-01-01T00:00:00Z",
  "last_active": "2024-01-01T00:00:00Z",
  "is_deleted": false,
  "deleted_at": null
}
```

### Remove Team Member
Soft deletes a team member. The member will be marked as deleted but not permanently removed from the team.

**Endpoint:** `DELETE /teams/{id}/members/{memberId}/`

**Path Parameters:**
- `id`: Team ID (string)
- `memberId`: Member ID (string)

**Response:** `204 No Content`

### Restore Team Member
Restores a previously deleted team member.

**Endpoint:** `POST /teams/{id}/members/{memberId}/restore/`

**Path Parameters:**
- `id`: Team ID (string)
- `memberId`: Member ID (string)

**Response:** `200 OK`

### Create Team Invitation
Creates an invitation for a new team member.

**Endpoint:** `POST /teams/{id}/invitations/`

**Path Parameters:**
- `id`: Team ID (string)

**Request Body:**
```json
{
  "email": "string",
  "role": "string"
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "email": "string",
  "role": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "expires_at": "2024-01-01T00:00:00Z",
  "is_used": false
}
```

### Accept Team Invitation
Accepts a team invitation.

**Endpoint:** `POST /teams/invitations/accept/`

**Request Body:**
```json
{
  "token": "string",
  "first_name": "string",
  "last_name": "string"
}
```

**Response:** `200 OK`

## Data Types

### Role
The role field can have the following values:
- `Admin User`
- `Account User`

### Team Object
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "members": [TeamMember],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "is_deleted": boolean,
  "deleted_at": "2024-01-01T00:00:00Z"
}
```

### Team Member Object
```json
{
  "id": "string",
  "user_id": "string",
  "role": "string",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "date_added": "2024-01-01T00:00:00Z",
  "last_active": "2024-01-01T00:00:00Z",
  "is_deleted": boolean,
  "deleted_at": "2024-01-01T00:00:00Z"
}
```

### Team Invitation Object
```json
{
  "id": "string",
  "team_id": "string",
  "email": "string",
  "role": "string",
  "token": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "expires_at": "2024-01-01T00:00:00Z",
  "is_used": boolean
}
```

### Validation Rules
- Team name is required
- Email addresses must be valid email format
- Role must be one of the defined role values
- First name and last name are required for team members
- Token is required for accepting invitations

### Timestamps
All timestamps are in ISO 8601 format with UTC timezone:
- Format: `YYYY-MM-DDThh:mm:ssZ`
- Example: `2024-01-01T00:00:00Z`

## Error Responses

### Common HTTP Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `204 No Content`: Request successful, no content to return
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
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