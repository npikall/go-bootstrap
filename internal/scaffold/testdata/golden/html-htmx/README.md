# demo

Scaffolded with [go-bootstrap](https://github.com/npikall/go-bootstrap).

## Development

```sh
make run
make watch   # live reload with air
make test
make lint
make fmt     # format with the formatters in .golangci.yml
make vuln    # check for known vulnerabilities with govulncheck
```

The application reads its configuration from environment variables only.
`make` exports the variables in `.env` to every recipe;
when running the binary directly, set them yourself.

Run `make help` to see all targets.
Missing development tools are installed system-wide on first use with
`go install` (pinned versions), after printing a warning.
