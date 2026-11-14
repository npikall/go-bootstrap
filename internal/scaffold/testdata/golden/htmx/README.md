# demo

Scaffolded with [go-bootstrap](https://github.com/npikall/go-bootstrap).

## Development

```sh
just run
just watch   # live reload with air
just test
just itest   # integration tests, requires Docker
just lint
just docker-run
```

The application reads its configuration from environment variables only.
`just` exports the variables in `.env` to every recipe;
when running the binary directly, set them yourself.

`just test` runs the tests with [gotestsum](https://github.com/gotestyourself/gotestsum),
`just test-watch` reruns them when files change.

## Database

Migrations in `database/migrations` are managed with [goose](https://github.com/pressly/goose),
configured through the `GOOSE_*` variables in `.env`:

```sh
just migrate-up       # apply pending migrations
just migrate-down     # roll back the latest migration
just migrate-status
just migrate-new add_email
```

[sqlc](https://sqlc.dev) generates type-safe Go code in `database/store` from the
queries in `database/queries` and the schema in `database/migrations`.
Run `just sqlc` after changing either, and use the result through
`database.New().Queries()`.

Run `just --list` to see all targets.
Missing development tools are installed system-wide on first use with
`go install` (pinned versions), after printing a warning.
