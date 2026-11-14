# go-bootstrap

Scaffold new Go web projects from a Jinja template, rendered with
[minijinja-go](https://github.com/mitsuhiko/minijinja/tree/main/minijinja-go).
It offers the same options as
[go-blueprint](https://github.com/Melkeydev/go-blueprint) and also lets you
choose the task runner.

## Install

```sh
go install github.com/npikall/go-bootstrap@latest
```

## Usage

Run `go-bootstrap create` in a terminal and answer the prompts. Options passed
as flags are not asked for again:

```sh
go-bootstrap create \
  --name myapp \
  --module github.com/me/myapp \
  --framework chi \
  --driver postgres \
  --feature htmx,tailwind,docker,githubaction \
  --tool sqlc,goose,gotestsum \
  --runner just \
  --git commit
```

| Flag               | Values                                                                                             | Default            |
| ------------------ | -------------------------------------------------------------------------------------------------- | ------------------ |
| `-name`, `-n`      | project name, also the output directory                                                            | required           |
| `-module`          | Go module path                                                                                     | project name       |
| `-framework`, `-f` | `chi`, `gin`, `fiber`, `gorilla/mux`, `httprouter`, `standard-library`, `echo`                     | `standard-library` |
| `-driver`, `-d`    | `mysql`, `postgres`, `sqlite`, `mongo`, `redis`, `scylla`, `none`                                  | `none`             |
| `-feature`         | `react`, `htmx`, `githubaction`, `websocket`, `tailwind`, `docker` (repeatable or comma separated) | none               |
| `-tool`            | `sqlc`, `goose`, `gotestsum` (repeatable or comma separated)                                       | none               |
| `-runner`          | `make`, `just`, `task`                                                                             | `make`             |
| `-git`, `-g`       | `commit`, `stage`, `skip`                                                                          | `commit`           |

`react` and `htmx` are mutually exclusive. `tailwind` without `react` adds a
templ frontend styled with the Tailwind standalone CLI.

`sqlc` and `goose` need a SQL driver (`mysql`, `postgres` or `sqlite`):

- `goose` adds `database/migrations` with a first migration, the `GOOSE_*`
  settings in `.env` and the `migrate-up`, `migrate-down`, `migrate-status` and
  `migrate-new` targets.
- `sqlc` adds `sqlc.yaml`, `database/queries` and the generated
  `database/store` package, exposed as `database.Service.Queries()`. It reads
  the schema from the goose migrations, or from `database/schema.sql` without
  goose. The `sqlc` target regenerates the code.
- `gotestsum` runs the `test` and `itest` targets and CI through gotestsum and
  adds a `test-watch` target.

After rendering, go-bootstrap runs `go mod tidy`, `templ generate` (when the
project uses templ) and sets up git.

## Generated projects

- Layout: `cmd/api`, `server`, `database`, `web` (templ) and `frontend` (React).
- Development tools (air, golangci-lint, templ, tailwindcss, sqlc, goose,
  gotestsum) are not in the `go.mod` tools directive, so their dependencies
  don't mix with the project's.
  Each runner target depends on an `install-<tool>` target. When the tool is not
  on `PATH`, it prints a warning and installs a pinned version system-wide with
  `go install` (into `go env GOBIN`, or `$(go env GOPATH)/bin`).
- A golangci-lint v2 config is included, and the generated code passes it.
- GitHub Actions workflows use pinned action SHAs and least-privilege
  permissions, and pass `zizmor`.
- Database integration tests use testcontainers behind the `integration` build
  tag (`make itest`), so `go test ./...` works without Docker.

## Template

The template lives in [`template/`](template) and follows the copier conventions
used by [py-bootstrap](https://github.com/npikall/py-bootstrap):

- Path segments are Jinja expressions. A segment that renders empty skips that
  file or directory, e.g. `{% if docker %}Dockerfile{% endif %}.jinja`.
- Files ending in `.jinja` are rendered and lose the suffix. Other files are
  copied verbatim.
- `template/_includes` holds partials used with `{% include %}`, such as the
  per-framework routes and per-driver database code.
- Templates render with `trim_blocks`, `lstrip_blocks` and strict undefined
  variables. The available variables are listed in
  [`internal/scaffold/options.go`](internal/scaffold/options.go).

Template paths become part of the Go module zip, so they may only use booleans
in conditions (`"` and `|` are invalid there). Go files in the template must end
in `.jinja`, otherwise the go tool treats them as packages of this module.

## Development

```sh
make test     # unit tests: render matrix, formatting checks, golden files
make golden   # update testdata/golden after changing the template
make e2e      # generate real projects, then build, vet, test and lint them (needs network)
make lint
```

`make test` renders every combination of framework, driver, runner and
representative feature sets. It checks that no file has consecutive or leading
blank lines, trailing whitespace, unrendered Jinja syntax or missing final
newlines, that Go files are gofmt-formatted, and that YAML/JSON files parse.
