# Product Requirements: CRUD To-do JSON REST API

## 1. Overview
Provide a simple JSON REST API to create, read, update, and delete To-do items. The API must be predictable, validate inputs strictly, and return consistent success and error response shapes.

## 2. Goals
- Full CRUD for To-do items
- Consistent JSON request/response formats
- Clear validation rules and error semantics
- Safe defaults (no silent coercion, no partial failures without explanation)

## 3. Non-Goals
- Authentication/authorization
- Multi-tenant workspaces/projects
- Search beyond basic filters

## 4. Conventions
### 4.1 Base URL
- `/api/v1`

### 4.2 Content Types
- Body requests require `Content-Type: application/json`
- Responses use `Content-Type: application/json`

### 4.3 Date/Time
- All timestamps are RFC3339 in UTC (e.g., `2026-05-04T12:34:56Z`)
- Server sets `createdAt` and `updatedAt`

### 4.4 Identifiers
- `id` is server-generated opaque string

## 5. Data Model
Field | Type | Notes
---|---|---
`id` | string | server-generated
`title` | string | required
`description` | string | optional
`completed` | boolean | default `false`
`dueAt` | string (RFC3339) \| null | optional/nullable
`createdAt` | string (RFC3339) | server-generated
`updatedAt` | string (RFC3339) | server-generated

## 6. Endpoints
### 6.1 Create
- `POST /api/v1/todos`
- Success: `201 Created`, `Location: /api/v1/todos/{id}`, body `{ "data": { ...todo... } }`

### 6.2 List
- `GET /api/v1/todos`
- Optional query params: `completed`, `limit` (1–100), `offset` (>=0), `sort` (`createdAt|updatedAt|dueAt|title`), `order` (`asc|desc`)
- Success: `200 OK`, body `{ "data": [...], "meta": { "limit": 50, "offset": 0, "count": N } }`

### 6.3 Get by ID
- `GET /api/v1/todos/{id}`
- Success: `200 OK`, body `{ "data": { ...todo... } }`

### 6.4 Update (Replace)
- `PUT /api/v1/todos/{id}` (idempotent)
- Success: `200 OK`, body `{ "data": { ...todo... } }`

### 6.5 Partial Update
- `PATCH /api/v1/todos/{id}`
- Success: `200 OK`, body `{ "data": { ...todo... } }`

### 6.6 Delete
- `DELETE /api/v1/todos/{id}`
- Success: `204 No Content` (empty body)

## 7. Validation Rules
- Reject invalid JSON bodies.
- Reject unknown top-level JSON fields (don’t ignore silently).
- `title`: trimmed length 1–200.
- `description`: trimmed length 0–2000.
- `dueAt`: RFC3339 datetime string or `null`.

Create (`POST`):
- `title` required; `completed` default `false`; others optional.

Replace (`PUT`):
- `title`, `description`, `completed`, `dueAt` all required (use `""` / `null` to clear).

Patch (`PATCH`):
- At least one mutable field present; omitted fields unchanged.

Path/query:
- validate `{id}` non-empty; validate query bounds and enum values.

## 8. Response Formats
### 8.1 Success
- Single: `{ "data": { ... } }`
- List: `{ "data": [ ... ], "meta": { "limit": 50, "offset": 0, "count": 1 } }`
- Delete: `204` no body

### 8.2 Errors
All non-2xx responses use:
```json
{
  "error": {
    "code": "validation_error",
    "message": "Validation failed",
    "details": [
      { "field": "title", "issue": "required" }
    ]
  }
}
```

## 9. Status Codes
- `200 OK`, `201 Created`, `204 No Content`
- `400 Bad Request` (invalid JSON / query / path)
- `404 Not Found` (missing id)
- `415 Unsupported Media Type` (missing/wrong content-type when body required)
- `422 Unprocessable Entity` (semantic validation errors; acceptable alternative is consistently using `400`)
- `500 Internal Server Error`

