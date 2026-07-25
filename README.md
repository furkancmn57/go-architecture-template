# Go Base Template (`graphql`)

GraphQL-only Go API template. Single binary entry at `src/main.go`.
In-memory storage — no GORM, Postgres, Redis, or REST.

| Doc | Path |
|-----|------|
| Architecture rules | [`.cursor/rules/go-artictech.mdc`](.cursor/rules/go-artictech.mdc) |
| Naming & folders | [`.cursor/STRUCTURE.md`](.cursor/STRUCTURE.md) |
| Commit messages | [`.cursor/COMMIT.md`](.cursor/COMMIT.md) |

## Stack

| Layer | Choice |
|-------|--------|
| HTTP host | Fiber v2 |
| API | GraphQL (`graphql-go`) at `/graphql` + GraphiQL |
| Storage | In-memory (process-local) |
| Validation | ozzo-validation |
| Logging | `log/slog` (text in local/dev, JSON elsewhere) |

REST + GORM + Postgres: [`main`](https://github.com/furkancmn57/go-architecture-template/tree/main).

## Layout

```text
src/
  main.go
  config/           # env loader + env/env.local, env/env.development
  extensions/       # GraphQL mount (/graphql + GraphiQL)
  common/           # apperr, logger
  constants/
  graphql/          # schema + resolvers → services
  services/         # {resource}/ + validations/
  models/           # requests/, responses/
```

## Run

Defaults to `APP_ENV=local` → loads `src/config/env/env.local`. No Docker deps.

```bash
make run        # APP_ENV=local
make run-dev    # APP_ENV=development
```

| Endpoint | URL |
|----------|-----|
| GraphQL | http://localhost:7090/graphql |
| Health | http://localhost:7090/health |

## Make

| Target | Description |
|--------|-------------|
| `make run` | Run API (`APP_ENV=local`) |
| `make run-dev` | Run API (`APP_ENV=development`) |
| `make build` | Build binary to `bin/` |
| `make clean` | Remove `bin/` |
| `make test` / `vet` | Tests and static analysis |
