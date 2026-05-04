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

