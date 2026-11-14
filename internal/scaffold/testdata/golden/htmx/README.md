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

Run `just --list` to see all targets.
Missing development tools are installed system-wide on first use with
`go install` (pinned versions), after printing a warning.
