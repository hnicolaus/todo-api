# todo-api (Go)

Simple in-memory CRUD To-do JSON REST API (Go standard library only).

## Run
```bash
cd todo-api
go run ./cmd/todo-api
```

Server listens on `:8080` by default. Override with:
```bash
ADDR=":9090" go run ./cmd/todo-api
```

## Endpoints
Base URL: `/api/v1`

- `GET /healthz`
- `POST /api/v1/todos`
- `GET /api/v1/todos`
- `GET /api/v1/todos/{id}`
- `PUT /api/v1/todos/{id}`
- `PATCH /api/v1/todos/{id}`
- `DELETE /api/v1/todos/{id}`

## Example curl
Create:
```bash
curl -i -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"Buy milk","description":"2 liters","dueAt":"2026-05-10T09:00:00Z"}'
```

List:
```bash
curl -s 'http://localhost:8080/api/v1/todos?limit=50&offset=0' | jq .
```

Get:
```bash
curl -s http://localhost:8080/api/v1/todos/<id> | jq .
```

Replace (PUT):
```bash
curl -s -X PUT http://localhost:8080/api/v1/todos/<id> \
  -H 'Content-Type: application/json' \
  -d '{"title":"Buy milk","description":"","completed":false,"dueAt":null}' | jq .
```

Patch:
```bash
curl -s -X PATCH http://localhost:8080/api/v1/todos/<id> \
  -H 'Content-Type: application/json' \
  -d '{"completed":true}' | jq .
```

Delete:
```bash
curl -i -X DELETE http://localhost:8080/api/v1/todos/<id>
```

## Tests
```bash
cd todo-api
go mod tidy
gofmt -w .
go test ./...
```

## Codex Prompt
```
Goal:
  Build a simple To-do REST API in Golang using an agentic workflow.

  Important:
  Use Codex subagents where useful. Spawn specialized subagents explicitly for PM, Technical Architect, and Software Engineers. Keep the main orchestrator
  focused on coordination, integration, and final validation.

  Workflow:

  1. PM Agent
  Ask a PM subagent to define the product requirements for a simple To-do app:
  - Create todo
  - Read/list todos
  - Update todo
  - Delete todo
  - Basic validation
  - JSON REST API
  - Expected success/error responses

  Output the PM result into:
  docs/PRODUCT_REQUIREMENTS.md

  2. Technical Architect Agent
  Ask a TA subagent to convert the PM requirements into a technical plan:
  - Go project structure
  - endpoint list
  - request/response schema
  - in-memory storage design
  - handler/service/repository boundaries
  - error handling convention
  - testing strategy
  - clear implementation tasks for Create, Read, Update, Delete

  Output the TA result into:
  docs/TECHNICAL_SPEC.md

  3. Software Engineer Subagents
  Spawn 4 Software Engineer subagents in parallel if possible:

  - SE-Create: implement Create Todo only
  - SE-Read: implement List/Get Todo only
  - SE-Update: implement Update Todo only
  - SE-Delete: implement Delete Todo only

  Each SE must:
  - Follow docs/TECHNICAL_SPEC.md
  - Avoid changing unrelated endpoint behavior
  - Add/update tests for their assigned area
  - Keep code simple and idiomatic Go

  4. Integration
  After the SE agents finish:
  - Review all changes
  - Resolve duplicated code, route conflicts, inconsistent response formats, or failing tests
  - Ensure the API is cohesive
  - Run formatting
  - Run all tests

  Commands to run:
  go mod tidy
  gofmt -w .
  go test ./...

  5. Final validation
  Create or update:
  README.md

  Include:
  - how to run the app
  - endpoint list
  - example curl commands
  - how to run tests

  Acceptance criteria:
  - `go test ./...` passes
  - CRUD endpoints work consistently
  - Code is simple and readable
  - README explains usage
  - No unnecessary dependencies
  - No database needed; use in-memory storage

  Implementation preference:
  Use only the Go standard library unless there is a strong reason not to.

  Before coding:
  First inspect the repository.
  If this is an empty directory, initialize a new Go module named `todo-api`.

  At the end, summarize:
  - agents used
  - files changed
  - tests run
  - any assumptions made
```
