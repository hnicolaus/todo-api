# Technical Spec: `todo-api` (Go stdlib only)

## 1. Scope
Implement a simple CRUD To-do REST API per `docs/PRODUCT_REQUIREMENTS.md`, using only the Go standard library and in-memory storage.

## 2. Project Layout
- `todo-api/cmd/todo-api/main.go`: entrypoint, server wiring
- `todo-api/internal/httpapi`: HTTP router, handlers, JSON helpers, error mapping
- `todo-api/internal/todo`: domain model, validation, service, repository interface + in-memory repo

## 3. HTTP API
### 3.1 Routes
Base: `/api/v1`
- `GET /healthz` -> `200` `{ "status": "ok" }`
- `POST /api/v1/todos` -> create
- `GET /api/v1/todos` -> list
- `GET /api/v1/todos/{id}` -> get
- `PUT /api/v1/todos/{id}` -> replace
- `PATCH /api/v1/todos/{id}` -> partial update
- `DELETE /api/v1/todos/{id}` -> delete

### 3.2 Request/Response Schemas
Success wrapper:
- Single: `{ "data": Todo }`
- List: `{ "data": []Todo, "meta": { "limit": number, "offset": number, "count": number } }`

Error wrapper:
```json
{
  "error": {
    "code": "validation_error",
    "message": "Validation failed",
    "details": [{ "field": "title", "issue": "required" }]
  }
}
```

`Todo` JSON fields:
- `id` (string), `title` (string), `description` (string), `completed` (bool)
- `dueAt` (string RFC3339 or null)
- `createdAt`/`updatedAt` (string RFC3339 UTC)

Body decoding rules:
- Reject unknown fields (use `json.Decoder.DisallowUnknownFields()`).
- Enforce `Content-Type: application/json` for requests with bodies.

## 4. In-memory Storage
Repository backed by `map[string]Todo` guarded by `sync.RWMutex`.
- ID generation: stdlib-only; use `crypto/rand` to generate random bytes and hex-encode to opaque id.
- List returns deterministic ordering: sort by a selected key (default `createdAt desc`).
- Support list filters/pagination: `completed`, `limit` (default 50, max 100), `offset` (default 0).

## 5. Boundaries
- `internal/todo`:
  - `Todo` model + validation helpers.
  - `Service` orchestrates validation and repo operations.
  - `Repository` interface for `Create/Get/List/Replace/Patch/Delete`.
- `internal/httpapi`:
  - `Server`/router wiring.
  - One handler file per operation to reduce merge conflicts.
  - Central JSON and error helpers for consistent responses.

## 6. Error Handling Convention
Define sentinel domain errors in `internal/todo`:
- `ErrNotFound`
- `ErrValidation`
- `ErrConflict` (reserved)

Map errors to HTTP in `internal/httpapi`:
- `ErrNotFound` -> `404` `{ error: { code: "not_found" } }`
- `ErrValidation` -> `422` `{ error: { code: "validation_error" } }`
- JSON parse/unknown fields -> `400` `{ error: { code: "invalid_json" } }`
- Wrong/missing content-type -> `415` `{ error: { code: "unsupported_media_type" } }`
- Default -> `500` `{ error: { code: "internal" } }`

## 7. Testing Strategy
Use `testing` + `net/http/httptest`.
- Handler tests: validate status codes, response shapes, and content-type handling.
- Repo tests: basic CRUD and list ordering/filtering.
- Consider `go test -race ./...` optionally; acceptance requires `go test ./...`.

## 8. Implementation Tasks (SE Ownership)
### 8.1 SE-Create
- Implement `POST /api/v1/todos` handler + service method.
- Validation: title required; default completed false; dueAt parse if provided.
- Tests: success `201`, bad content-type `415`, invalid JSON `400`, validation error `422`.

### 8.2 SE-Read
- Implement `GET /api/v1/todos` (list) and `GET /api/v1/todos/{id}` (get).
- Query parsing: `completed`, `limit`, `offset`, `sort`, `order`.
- Tests: list empty, list filtered, get not-found -> `404`.

### 8.3 SE-Update
- Implement `PUT` (replace) and `PATCH` (partial update).
- Ensure `updatedAt` changes on success; `PUT` requires all fields; `PATCH` requires at least one field.
- Tests: replace success, patch success, not-found -> `404`, validation -> `422`, unknown field -> `400`.

### 8.4 SE-Delete
- Implement `DELETE /api/v1/todos/{id}` returning `204`.
- Missing resource -> `404`.
- Tests: delete success then get -> `404`, delete not-found -> `404`.

