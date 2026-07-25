# Naming & Folder Conventions (`graphql`)

GraphQL-only branch. In-memory storage. No GORM / REST / OpenAPI.

## Models — one type per file

| Type | Path |
|------|------|
| `CreateTodoRequest` | `src/models/requests/create_todo_request.go` |
| `UpdateTodoRequest` | `src/models/requests/update_todo_request.go` |
| `TodoResponse` | `src/models/responses/todo_response.go` |

```text
src/models/requests/{action}_{resource}_request.go
src/models/responses/{resource}_response.go
```

## Validations

```text
src/services/{resource}/validations/
  create_{resource}_request.go
  update_{resource}_request.go
```

Called at service entry: `validations.CreateTodoRequest(req)`.

## Error codes

```text
src/constants/errors.go
src/constants/{domain}_errors.go
```

## GraphQL

```text
src/graphql/schema.go
src/graphql/{resource}.go
src/extensions/graphql.go   # /graphql + GraphiQL
```

Resolvers → `services/` only.

## Service

```text
src/services/{resource}/
  service.go        # use-cases + in-memory map
  validations/
```

## New resource checklist

1. `services/{resource}/` — service + validations
2. `models/` DTOs
3. `graphql/` types + fields
4. Wire in `main.go` → `NewSchema(...)`
