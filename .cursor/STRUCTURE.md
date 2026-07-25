# Naming & Folder Conventions

This project follows the NotificationApi horizontal layout. Use this file as
the naming source of truth when adding models and folders.

## Schema — code-first + Fluent mappings (NotificationApi-style)

| Layer | Role |
|-------|------|
| `entities/` | Clean POCO (fields only). No column/type GORM tags. |
| `mappings/{name}_map.go` | `TableName`, `Entity()`, Fluent `Columns()`, `Migrate(tx)`. |
| `mappings/fluent.go` | `Column` + `ApplyColumns` (IsRequired / HasMaxLength equivalent). |
| `migrations/` | Versioned steps call `mappings.XMap{}.Migrate(tx)`. |
| `schema_migrations` | Applied version history. |

```text
src/data/entities/todo.go                # POCO (TableName matches map)
src/data/mappings/todo_map.go            # TodoMap: Entity + Columns + Migrate
src/data/mappings/fluent.go              # ApplyColumns helper
src/data/migrations/000001_create_todos.go
src/data/migrate.go                      # Setup / Migrate runner
src/data/ensure.go                       # CREATE DATABASE if missing
```

Startup: `EnsureDatabase` → connect → `Migrate` (pending only).

`Migrate` = `AutoMigrate(entity)` then `ApplyColumns` for Fluent rules. No shadow
column structs. Entity must not import `mappings` (import cycle).

**New table checklist**

1. `entities/{name}.go` — fields only + `TableName()` matching the map
2. `mappings/{name}_map.go` — `Entity`, `Columns`, `Migrate` (`AutoMigrate` + `ApplyColumns`)
3. `migrations/00000N_....go` — `Up: mappings.XMap{}.Migrate`
4. Restart — row in `schema_migrations`

Do not put domain schema tags on entities. Do not call `AutoMigrate` from services/controllers.

## Models — one type per file

Each request and response DTO lives in its own file. File name = snake_case of
the type name.

| Type | Path |
|------|------|
| `CreateTodoRequest` | `src/models/requests/create_todo_request.go` |
| `UpdateTodoRequest` | `src/models/requests/update_todo_request.go` |
| `TodoResponse` | `src/models/responses/todo_response.go` |

**Pattern**

```text
src/models/requests/{action}_{resource}_request.go   → type {Action}{Resource}Request
src/models/responses/{resource}_response.go          → type {Resource}Response
```

Do not put multiple request/response structs in the same file.

## Validations — folder under the service

```text
src/services/{resource}/validations/
  create_{resource}_request.go   → validations.Create{Resource}Request
  update_{resource}_request.go   → validations.Update{Resource}Request
```

Called from the service at method entry: `validations.CreateTodoRequest(req)`.

## Error codes

API JSON `code` values live in `constants/`:

```text
src/constants/errors.go           # shared
src/constants/{domain}_errors.go  # e.g. TODO_NOT_FOUND
```

Call sites: `apperr.NotFound(constants.TodoNotFound, "todo not found")`.

## HTTP — controllers

```text
src/controllers/v1/{resource}.go   # handlers + Register(api)
```

`main` → `api := app.Group("/api/v1")` → `NewXController(svc).Register(api)`.

## Quick checklist for a new resource

1. Entity (POCO) + Fluent map (`Entity`/`Columns`/`Migrate`) + versioned migration
2. Requests / responses (one type per file)
3. Service + validations
4. Controller handlers + `Register(api)` + swag + `make openapi`
5. Wire in `main.go`

## GraphQL

Not on `main`. Use branch
[`feature/graphql`](https://github.com/furkancmn57/go-base-template/tree/feature/graphql)
(`git checkout feature/graphql`).
