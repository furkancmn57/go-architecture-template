# Go Base Template (`feature/graphql`)

Lean GraphQL-only template. In-memory storage — no GORM, Postgres, Redis, or REST.

| Doc | Path |
|-----|------|
| Rules | [`.cursor/rules/go-artictech.mdc`](.cursor/rules/go-artictech.mdc) |
| Structure | [`.cursor/STRUCTURE.md`](.cursor/STRUCTURE.md) |
| Commits | [`.cursor/COMMIT.md`](.cursor/COMMIT.md) |

## Stack

Fiber v2 · graphql-go · ozzo-validation · slog · in-memory store

REST + GORM: [`main`](https://github.com/furkancmn57/go-base-template/tree/main)

## Layout

```text
src/
  main.go
  config/           # APP_ENV, APP_PORT
  extensions/       # GraphQL mount
  common/           # apperr, logger
  constants/
  graphql/          # schema + resolvers → services
  services/todo/    # service + validations
  models/           # requests/, responses/
```

## Run

```bash
make run        # http://localhost:7090/graphql
make run-dev
make clean      # rm -rf bin/
```

| Endpoint | URL |
|----------|-----|
| GraphQL | http://localhost:7090/graphql |
| Health | http://localhost:7090/health |
