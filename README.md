# Go Base Template (`feature/graphql`)

Same stack as `main`, plus optional GraphQL transport at `/graphql`.

| Doc | Path |
|-----|------|
| Architecture rules | [`.cursor/rules/go-artictech.mdc`](.cursor/rules/go-artictech.mdc) |
| Naming & folders | [`.cursor/STRUCTURE.md`](.cursor/STRUCTURE.md) |
| Commit messages | [`.cursor/COMMIT.md`](.cursor/COMMIT.md) |

## Stack

| Layer | Choice |
|-------|--------|
| HTTP | Fiber v2 |
| DB | GORM + Postgres (code-first Fluent maps + versioned migrations) |
| Cache | Redis |
| Validation | ozzo-validation |
| Logging | `log/slog` (text in local/dev, JSON elsewhere) |
| API docs | OpenAPI (swaggo) at `/openapi/*` |
| GraphQL | Optional (`GRAPHQL_ENABLED=true`) — same services as REST |

## Layout

```text
src/
  main.go
  config/           # env loader + env/env.local, env/env.development
  extensions/       # DB, Redis, health, OpenAPI, GraphQL
  common/           # Model, WriteJSON, apperr, logger
  constants/
  interfaces/       # Cache
  controllers/v1/   # thin handlers + Register(api)
  graphql/          # schema/resolvers → same services
  services/         # {resource}/ + validations/; cache/
  data/             # entities (POCO) + Fluent mappings + versioned migrations
  models/           # requests/, responses/
  docs/             # swag output (make openapi)
```

## Run

Defaults to `APP_ENV=local` → loads `src/config/env/env.local`.

```bash
docker compose up -d postgres redis
make run        # APP_ENV=local
make run-dev    # APP_ENV=development
```

Enable GraphQL:

```bash
GRAPHQL_ENABLED=true make run
```

| Endpoint | URL |
|----------|-----|
| API | http://localhost:7090/api/v1 |
| Health | http://localhost:7090/health |
| OpenAPI | http://localhost:7090/openapi/index.html |
| GraphQL | http://localhost:7090/graphql |

REST-only template: [`main`](https://github.com/furkancmn57/go-base-template/tree/main).

## Make

| Target | Description |
|--------|-------------|
| `make run` | Run API (`APP_ENV=local`) |
| `make run-dev` | Run API (`APP_ENV=development`) |
| `make build` | Build binary to `bin/` |
| `make openapi` | Regenerate `src/docs` from controller annotations |
| `make docker-up` / `docker-down` | Local Postgres, Redis |
| `make test` / `vet` | Tests and static analysis |

```bash
go install github.com/swaggo/swag/cmd/swag@latest   # make openapi
```
