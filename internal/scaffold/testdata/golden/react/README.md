# demo

Scaffolded with [go-bootstrap](https://github.com/npikall/go-bootstrap).

## Development

```sh
task run     # API on :8080 and the frontend dev server
task watch   # live reload with air
task test
task itest   # integration tests, requires Docker
task lint
task docker-run
```

The application reads its configuration from environment variables only.
`task` exports the variables in `.env` to every task;
when running the binary directly, set them yourself.

`task test` runs the tests with [gotestsum](https://github.com/gotestyourself/gotestsum),
`task test-watch` reruns them when files change.

Run `task --list` to see all targets.
Missing development tools are installed system-wide on first use with
`go install` (pinned versions), after printing a warning.
