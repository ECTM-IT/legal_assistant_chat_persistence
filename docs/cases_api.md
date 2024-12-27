# Cases API Documentation

## Overview
The Cases API provides endpoints for managing legal cases, including case creation, updates, document management, collaborator management, and message feedback. Each case can contain messages, documents, and collaborators.

## Base URL
All endpoints are relative to the base URL of your API server.

## Authentication
All endpoints require authentication. Include your authentication token in the Authorization header:
```
Authorization: Bearer <your_token>
```

## Endpoints

### Get All Cases
Retrieves all cases accessible to the authenticated user.

**Endpoint:** `GET /api-cases/`

**Response:** `200 OK`
```json
[
  {
    "id": "string",
    "name": "string",
    "creator_id": "string",
    "messages": [Message],
    "collaborators": [Collaborator],
    "skills": ["string"],
    "documents": [Document],
    "action": "string",
    "agent_id": "string",
    "creation_date": "2024-01-01T00:00:00Z",
    "last_edit": "2024-01-01T00:00:00Z",
    "share": boolean,
    "is_archived": boolean
  }
]
```

### Get Cases by Creator
Retrieves all cases created by the authenticated user.

**Endpoint:** `GET /cases-user/`

**Response:** `200 OK`
```json
[
  {
    "id": "string",
    "name": "string",
    "creator_id": "string",
    "messages": [Message],
    "collaborators": [Collaborator],
    "skills": ["string"],
    "documents": [Document],
    "action": "string",
    "agent_id": "string",
    "creation_date": "2024-01-01T00:00:00Z",
    "last_edit": "2024-01-01T00:00:00Z",
    "share": boolean,
    "is_archived": boolean
  }
]
```

### Get Case by ID
Retrieves a specific case by its ID.

**Endpoint:** `GET /cases/{id}/`

**Path Parameters:**
- `id`: Case ID (string)

**Response:** `200 OK`
```json
{
  "id": "string",
  "name": "string",
  "creator_id": "string",
  "messages": [Message],
  "collaborators": [Collaborator],
  "skills": ["string"],
  "documents": [Document],
  "action": "string",
  "agent_id": "string",
  "creation_date": "2024-01-01T00:00:00Z",
  "last_edit": "2024-01-01T00:00:00Z",
  "share": boolean,
  "is_archived": boolean
}
```

### Create Case
Creates a new case.

**Endpoint:** `POST /cases-create/`

**Request Body:**
```json
{
  "name": "string",
  "skills": ["string"],
  "action": "string",
  "agent_id": "string",
  "share": boolean
}
```

**Response:** `201 Created`
```json
{
  "id": "string",
  "name": "string",
  "creator_id": "string",
  "messages": [],
  "collaborators": [],
  "skills": ["string"],
  "documents": [],
  "action": "string",
  "agent_id": "string",
  "creation_date": "2024-01-01T00:00:00Z",
  "last_edit": "2024-01-01T00:00:00Z",
  "share": boolean,
  "is_archived": false
}
```

### Update Case
Updates an existing case.

**Endpoint:** `PATCH /cases/{id}/`

**Path Parameters:**
- `id`: Case ID (string)

**Request Body:**
```json
{
  "name": "string",
  "skills": ["string"],
  "action": "string",
  "share": boolean,
  "is_archived": boolean
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "name": "string",
  "creator_id": "string",
  "messages": [Message],
  "collaborators": [Collaborator],
  "skills": ["string"],
  "documents": [Document],
  "action": "string",
  "agent_id": "string",
  "creation_date": "2024-01-01T00:00:00Z",
  "last_edit": "2024-01-01T00:00:00Z",
  "share": boolean,
  "is_archived": boolean
}
```

### Delete Case
Deletes a case.

**Endpoint:** `DELETE /cases/{id}/`

**Path Parameters:**
- `id`: Case ID (string)

**Response:** `204 No Content`

### Add Collaborator to Case
Adds a collaborator to a case.

**Endpoint:** `POST /case-add-user/{id}/`

**Path Parameters:**
- `id`: Case ID (string)

**Request Body:**
```json
{
  "id": "string",
  "edit": boolean
}
```

**Response:** `200 OK`

### Remove Collaborator from Case
Removes a collaborator from a case.

**Endpoint:** `DELETE /case-remove-user/{id}/{userID}/`

**Path Parameters:**
- `id`: Case ID (string)
- `userID`: User ID (string)

**Response:** `204 No Content`

### Add Document to Case
Adds a document to a case.

**Endpoint:** `POST /case-add-document/{id}/`

**Path Parameters:**
- `id`: Case ID (string)

**Request Body:**
```json
{
  "sender": "string",
  "file_name": "string",
  "file_type": "string",
  "file_content": "string",
  "collaborators": [
    {
      "email": "string",
      "edit": boolean
    }
  ]
}
```

**Response:** `200 OK`
```json
{
  "id": "string",
  "created_by": "string",
  "sender": "string",
  "file_name": "string",
  "file_type": "string",
  "file_content": "string",
  "collaborators": [DocumentCollaborator],
  "upload_date": "2024-01-01T00:00:00Z",
  "modified_date": "2024-01-01T00:00:00Z"
}
```

### Update Document
Updates a document in a case.

**Endpoint:** `PATCH /case-update-document/{id}/document/{documentID}/`

**Path Parameters:**
- `id`: Case ID (string)
- `documentID`: Document ID (string)

**Request Body:**
```json
{
  "file_content": "string",
  "file_name": "string"
}
```

**Response:** `200 OK`

### Add Document Collaborator
Adds a collaborator to a document.

**Endpoint:** `POST /case-add-document-collaborator/{id}/document/{documentID}/`

**Path Parameters:**
- `id`: Case ID (string)
- `documentID`: Document ID (string)

**Request Body:**
```json
{
  "email": "string",
  "edit": boolean
}
```

**Response:** `200 OK`

### Delete Document
Deletes a document from a case.

**Endpoint:** `DELETE /case-remove-document/{id}/{documentID}/`

**Path Parameters:**
- `id`: Case ID (string)
- `documentID`: Document ID (string)

**Response:** `204 No Content`

### Add Feedback to Message
Adds feedback to a message in a case.

**Endpoint:** `POST /case-add-feedback-to-message/{id}/{messageId}/`

**Path Parameters:**
- `id`: Case ID (string)
- `messageId`: Message ID (string)

**Request Body:**
```json
{
  "score": "string",
  "reasons": ["string"],
  "comment": "string"
}
```

**Response:** `200 OK`

### Get Message Feedback
Gets feedback for a specific message.

**Endpoint:** `GET /case-get-feedback-from-message/{id}/{creatorId}/{messageId}/`

**Path Parameters:**
- `id`: Case ID (string)
- `creatorId`: Creator ID (string)
- `messageId`: Message ID (string)

**Response:** `200 OK`
```json
{
  "id": "string",
  "case_id": "string",
  "message_id": "string",
  "creator_id": "string",
  "score": "string",
  "reasons": ["string"],
  "comment": "string",
  "creation_date": "2024-01-01T00:00:00Z"
}
```

## Data Types

### Case Object
```json
{
  "id": "string",
  "name": "string",
  "creator_id": "string",
  "messages": [Message],
  "collaborators": [Collaborator],
  "skills": ["string"],
  "documents": [Document],
  "action": "string",
  "agent_id": "string",
  "creation_date": "2024-01-01T00:00:00Z",
  "last_edit": "2024-01-01T00:00:00Z",
  "share": boolean,
  "is_archived": boolean
}
```

### Message Object
```json
{
  "id": "string",
  "sender": "string",
  "recipient": "string",
  "content": "string",
  "document_path": "string",
  "function_call": boolean,
  "feedbacks": [Feedback],
  "skills": ["string"],
  "agent": "string"
}
```

### Collaborator Object
```json
{
  "id": "string",
  "edit": boolean
}
```

### Document Object
```json
{
  "id": "string",
  "created_by": "string",
  "sender": "string",
  "file_name": "string",
  "file_type": "string",
  "file_content": "string",
  "collaborators": [DocumentCollaborator],
  "upload_date": "2024-01-01T00:00:00Z",
  "modified_date": "2024-01-01T00:00:00Z"
}
```

### DocumentCollaborator Object
```json
{
  "email": "string",
  "edit": boolean
}
```

### Feedback Object
```json
{
  "id": "string",
  "case_id": "string",
  "message_id": "string",
  "creator_id": "string",
  "score": "string",
  "reasons": ["string"],
  "comment": "string",
  "creation_date": "2024-01-01T00:00:00Z"
}
```

### Validation Rules
- Case name is required
- Creator ID is required
- Agent ID is required when creating a case
- File name and type are required for documents
- Email must be valid format for document collaborators
- Message ID is required for feedback
- Score is required for feedback

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