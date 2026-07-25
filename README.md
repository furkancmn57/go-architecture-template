# Go Base Template

Horizontal-layer Go API template. Single binary entry at `src/main.go`.

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

## Layout

```text
src/
  main.go
  config/           # env loader + env/env.local, env/env.development
  extensions/       # DB, Redis, health, OpenAPI
  common/           # Model, WriteJSON, apperr, logger
  constants/
  interfaces/       # Cache
  controllers/v1/   # thin handlers + Register(api)
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

| Endpoint | URL |
|----------|-----|
| API | http://localhost:7090/api/v1 |
| Health | http://localhost:7090/health |
| OpenAPI | http://localhost:7090/openapi/index.html |

## GraphQL (optional)

`main` is REST-only. GraphQL lives on a separate branch:

| | |
|--|--|
| Branch | [`graphql`](https://github.com/furkancmn57/go-base-template/tree/graphql) |
| Checkout | `git fetch origin && git checkout graphql` |

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
