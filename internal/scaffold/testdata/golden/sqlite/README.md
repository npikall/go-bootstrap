# demo

Scaffolded with [go-bootstrap](https://github.com/npikall/go-bootstrap).

## Development

```sh
make run
make watch   # live reload with air
make test
make lint
```

`make test` runs the tests with [gotestsum](https://github.com/gotestyourself/gotestsum),
`make test-watch` reruns them when files change.

## Database

Migrations in `database/migrations` are managed with [goose](https://github.com/pressly/goose),
configured through the `GOOSE_*` variables in `.env`:

```sh
make migrate-up       # apply pending migrations
make migrate-down     # roll back the latest migration
make migrate-status
make migrate-new name=add_email
```

[sqlc](https://sqlc.dev) generates type-safe Go code in `database/store` from the
queries in `database/queries` and the schema in `database/migrations`.
Run `make sqlc` after changing either, and use the result through
`database.New().Queries()`.

Run `make help` to see all targets.
Missing development tools are installed system-wide on first use with
`go install` (pinned versions), after printing a warning.
