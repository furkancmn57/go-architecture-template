# Commit Convention (Conventional Commits)

All commits in this repo follow [Conventional Commits](https://www.conventionalcommits.org/).
Agents and humans must use this format when creating git commits.

## Format

```text
<type>(<optional-scope>): <short description>

[optional body]

[optional footer]
```

- Subject line: imperative mood, lowercase after the colon, no trailing period, ≤72 chars preferred
- Body: why (not what), wrap ~72 chars; separate from subject with a blank line
- One logical change per commit when practical

## Types

| Type | When |
|------|------|
| `feat` | New user-facing capability or API |
| `fix` | Bug fix |
| `docs` | Documentation / rules / README only |
| `style` | Formatting; no logic change |
| `refactor` | Code change that is not feat/fix |
| `perf` | Performance improvement |
| `test` | Tests only |
| `build` | Go modules, Makefile, Docker, deps |
| `ci` | CI configuration |
| `chore` | Maintenance that does not fit above |
| `revert` | Revert a previous commit |

## Scopes (optional, prefer when clear)

Examples used in this project:

- `todo` — todo service / store / resolvers
- `graphql` — schema / transport
- `config` — env / config packages
- `cursor` — `.cursor/` rules and project guides

Omit scope when the change is cross-cutting.

## Examples

```text
feat(todo): add complete endpoint

docs(cursor): move STRUCTURE naming guide under .cursor

build: add graphql-go dependencies

chore: ignore .idea IDE files
```

## Forbidden

- Vague subjects: `update`, `fix stuff`, `wip`, `changes`
- Commit messages that only list file names
- Secrets / `.env` credentials in any commit
- `--no-verify` / skipping hooks unless explicitly requested
```

