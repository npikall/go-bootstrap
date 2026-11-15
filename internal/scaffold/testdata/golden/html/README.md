# demo

Scaffolded with [go-bootstrap](https://github.com/npikall/go-bootstrap).

## Development

```sh
task run
task watch   # live reload with air
task test
task lint
task fmt     # format with the formatters in .golangci.yml
task vuln    # check for known vulnerabilities with govulncheck
```

The application reads its configuration from environment variables only.
`task` exports the variables in `.env` to every task;
when running the binary directly, set them yourself.

Run `task --list` to see all targets.
Missing development tools are installed system-wide on first use with
`go install` (pinned versions), after printing a warning.
