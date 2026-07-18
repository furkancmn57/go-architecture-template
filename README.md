# Go Base Template

Horizontal-layer Go API template (NotificationApi-style). Single binary entry at `src/main.go`.

| Doc | Path |
|-----|------|
| Architecture rules | [`.cursor/rules/go-artictech.mdc`](.cursor/rules/go-artictech.mdc) |
| Naming & folders | [`.cursor/STRUCTURE.md`](.cursor/STRUCTURE.md) |
| Commit messages | [`.cursor/COMMIT.md`](.cursor/COMMIT.md) |

## Stack

| Layer | Choice |
|-------|--------|
| HTTP | Fiber v2 |
| DB | GORM + Postgres |
| Cache | Redis |
| Messaging | RabbitMQ |
| Validation | ozzo-validation |
| API docs | OpenAPI (swaggo) at `/openapi/*` |

## Layout

```text
src/
  main.go
  config/           # env loader
  extensions/       # DB, Redis, RabbitMQ, health, OpenAPI
  common/           # Model, WriteJSON, BaseConsumer, apperr
  constants/
  interfaces/       # Publisher, Subscriber, Cache
  controllers/v1/   # thin HTTP handlers
  services/         # {resource}/ + validations/; cache/
  data/             # entities, mappings, postgres, migrate
  models/           # requests/, responses/
  events/           # {domain}/
  consumers/        # {domain}/
  docs/             # swag output (make openapi)
```

## Run

```bash
cp src/config/env/.env.example .env
docker compose up -d postgres redis rabbitmq
make run
```

| Endpoint | URL |
|----------|-----|
| API | http://localhost:8080/api/v1 |
| Health | http://localhost:8080/health |
| OpenAPI | http://localhost:8080/openapi/index.html |

## GraphQL (optional)

`main` REST-only. GraphQL transport lives on a separate branch:

| | |
|--|--|
| Branch | [`feature/graphql`](https://github.com/furkancmn57/go-base-template/tree/feature/graphql) |
| Checkout | `git fetch origin && git checkout feature/graphql` |

## Make

| Target | Description |
|--------|-------------|
| `make run` | Run the API |
| `make build` | Build binary to `bin/` |
| `make openapi` | Regenerate `src/docs` from controller annotations |
| `make docker-up` / `docker-down` | Local Postgres, Redis, RabbitMQ |
| `make test` / `vet` | Tests and static analysis |

`make openapi` requires: `go install github.com/swaggo/swag/cmd/swag@latest`
