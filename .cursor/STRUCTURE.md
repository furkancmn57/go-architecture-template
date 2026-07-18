# Naming & Folder Conventions

This project follows the NotificationApi horizontal layout. Use this file as
the naming source of truth when adding models, events, and consumers.

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

Examples (NotificationApi → Go):

- `TemplateCreateRequest.cs` → `create_template_request.go` / `CreateTemplateRequest`
- `TemplateFilterRequest.cs` → `filter_template_request.go` / `FilterTemplateRequest`

Do not put multiple request/response structs in the same file.

## Validations — folder under the service

```text
src/services/{resource}/validations/
  create_{resource}_request.go   → validations.Create{Resource}Request
  update_{resource}_request.go   → validations.Update{Resource}Request
```

Called from the service at method entry: `validations.CreateTodoRequest(req)`.

## RabbitMQ

```text
src/interfaces/publisher.go    # Publisher port
src/interfaces/subscriber.go   # Subscriber + Handler ports
src/extensions/rabbitmq.go     # AddRabbitMQ + Publish/Subscribe/Close
```

Exchange name: `constants.RabbitMQExchange`. Topic strings: `constants/{domain}_topics.go`.

## Error codes

API JSON `code` values live in `constants/` (not hardcoded in services/controllers):

```text
src/constants/errors.go           # shared: INVALID_REQUEST_BODY, VALIDATION_ERROR, …
src/constants/{domain}_errors.go  # e.g. todo_errors.go → TODO_NOT_FOUND
```

Call sites: `apperr.NotFound(constants.TodoNotFound, "todo not found")`.

## Events — domain folder, one event per file

Folder = domain / entity the event belongs to (NotificationApi `Events/News/`).

```text
src/events/{domain}/
  {action}_event.go         # one payload struct per file

src/constants/{domain}_topics.go   # topic string constants
```

Todo example:

```text
src/events/todo/
  created_event.go      → CreatedEvent
  updated_event.go      → UpdatedEvent
  completed_event.go    → CompletedEvent
  deleted_event.go      → DeletedEvent

src/constants/todo_topics.go
  TodoCreated / TodoUpdated / TodoCompleted / TodoDeleted
```

Topic naming: `{module}.{action}` (e.g. `todo.completed`). Topic strings live in
`constants/`, payload structs live in `events/`.

## Consumers — domain folder + {Verb}When{Action}

Folder = domain / entity the consumer works on (NotificationApi `Consumers/News/`).

Naming mirrors NotificationApi (`SendMailWhenForgotPassword`):

```text
src/consumers/{domain}/{verb}_when_{action}.go
  → type {Verb}When{Action}
  → New{Verb}When{Action}(...)
  → Register(subscriber)
```

Todo example:

```text
src/consumers/todo/log_when_todo_completed.go
  → LogWhenTodoCompleted
```

Rules:

1. One consumer type per file.
2. Folder name = the entity/domain the consumer handles (`todo`, …).
3. Type name = verb + `When` + action/event (e.g. `LogWhenTodoCompleted`, `SendMailWhenForgotPassword`).
4. Wire each consumer explicitly in `main.go` via `Register`.

## Quick checklist for a new resource

1. Entity: `data/entities/{name}.go` + `data/mappings/{name}.go` + `data/migrate.go`
2. Requests: one file each under `models/requests/`
3. Responses: one file each under `models/responses/`
4. Service: `services/{name}/service.go` + `validations/` (one file per request);
   shared helper `apperr.FromValidation`
5. Events (if any): `events/{name}/{action}_event.go` + topic consts in `constants/{name}_topics.go`
6. Consumers (if any): `consumers/{name}/{verb}_when_{action}.go`
7. Controller: `controllers/v1/{name}.go` + `Register(api)` for its routes
8. Controller annotations (`@Summary`, `@Router`, …) + `make openapi` (commit `src/docs`)
9. Wire in `main.go`: `NewXController(svc).Register(api)`

## GraphQL

Not on `main`. Use branch
[`feature/graphql`](https://github.com/furkancmn57/go-base-template/tree/feature/graphql)
(`git checkout feature/graphql`).
