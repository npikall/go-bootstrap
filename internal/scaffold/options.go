package scaffold

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/mod/module"
)

var (
	Frameworks = []string{"chi", "gin", "fiber", "gorilla/mux", "httprouter", "standard-library", "echo"}
	Databases  = []string{"mysql", "postgres", "sqlite", "mongo", "redis", "scylla", "none"}
	Features   = []string{"react", "htmx", "githubaction", "websocket", "tailwind", "docker"}
	Tools      = []string{"sqlc", "goose", "gotestsum"}
	Runners    = []string{"make", "just", "task"}
	GitModes   = []string{"commit", "stage", "skip"}
)

// Tool versions pinned into the generated task runner files.
var toolVersions = map[string]string{
	"air_version":           "v1.67.4",
	"golangci_lint_version": "v2.15.0",
	"tailwind_version":      "v4.3.5",
	"sqlc_version":          "v1.31.1",
	"goose_version":         "v3.30.1",
	"gotestsum_version":     "v1.13.0",
}

// sqlDialects maps each SQL database to the names sqlc and goose use for it.
// Only these databases support the tools that need SQL.
var sqlDialects = map[string]struct{ sqlcEngine, gooseDriver, param, idType string }{
	"postgres": {"postgresql", "postgres", "$1", "BIGSERIAL"},
	"mysql":    {"mysql", "mysql", "?", "BIGINT AUTO_INCREMENT"},
	"sqlite":   {"sqlite", "sqlite3", "?", "INTEGER"},
}

// sqlTools are the tools that only work with a SQL database.
var sqlTools = []string{"sqlc", "goose"}

var namePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)

// Options describes the project to generate.
type Options struct {
	Name       string
	ModulePath string
	GoVersion  string
	Framework  string
	Database   string
	Features   []string
	Tools      []string
	Runner     string
	Git        string
}

// Defaults returns the options used when nothing else is specified.
func Defaults() Options {
	return Options{
		Framework: "standard-library",
		Database:  "none",
		Runner:    "make",
		Git:       "commit",
	}
}

// Validate reports every invalid option at once.
func (o Options) Validate() error {
	var errGoVersion error
	if o.GoVersion == "" {
		errGoVersion = errors.New("go version must not be empty")
	}
	return errors.Join(ValidateName(o.Name), ValidateModulePath(o.ModulePath), errGoVersion, o.ValidateChoices())
}

// ValidateName checks that name is usable as directory and default module path.
func ValidateName(name string) error {
	if !namePattern.MatchString(name) {
		return fmt.Errorf("invalid project name %q: use letters, digits, '.', '_' or '-'", name)
	}
	return nil
}

// ValidateModulePath checks a module path, the empty string means "use the project name".
func ValidateModulePath(path string) error {
	if path == "" {
		return nil
	}
	if err := module.CheckPath(path); err != nil {
		return fmt.Errorf("invalid module path: %w", err)
	}
	return nil
}

// ValidateChoices checks the enumerated options and feature combinations.
func (o Options) ValidateChoices() error {
	return errors.Join(
		oneOf("framework", o.Framework, Frameworks),
		oneOf("database driver", o.Database, Databases),
		oneOf("runner", o.Runner, Runners),
		oneOf("git option", o.Git, GitModes),
		ValidateFeatures(o.Features),
		ValidateTools(o.Tools, o.Database),
	)
}

// ValidateTools checks tool names and that tools needing SQL get a SQL database.
func ValidateTools(tools []string, database string) error {
	var errs []error
	_, isSQL := sqlDialects[database]
	for _, tool := range tools {
		errs = append(errs, oneOf("tool", tool, Tools))
		if slices.Contains(sqlTools, tool) && !isSQL {
			errs = append(errs, fmt.Errorf("tool %s requires a SQL database driver (mysql, postgres or sqlite), got %q", tool, database))
		}
	}
	return errors.Join(errs...)
}

// ValidateFeatures checks feature names and combinations.
func ValidateFeatures(features []string) error {
	var errs []error
	for _, f := range features {
		errs = append(errs, oneOf("feature", f, Features))
	}
	if slices.Contains(features, "react") && slices.Contains(features, "htmx") {
		errs = append(errs, errors.New("features react and htmx are mutually exclusive"))
	}
	return errors.Join(errs...)
}

// Has reports whether the feature is enabled.
func (o Options) Has(feature string) bool {
	return slices.Contains(o.Features, feature)
}

// Uses reports whether the tool is enabled.
func (o Options) Uses(tool string) bool {
	return slices.Contains(o.Tools, tool)
}

func (o Options) module() string {
	if o.ModulePath != "" {
		return o.ModulePath
	}
	return o.Name
}

// usesTempl reports whether the project gets a templ based web frontend.
func (o Options) usesTempl() bool {
	return (o.Has("htmx") || o.Has("tailwind")) && !o.Has("react")
}

// Data returns the variables available inside templates.
func (o Options) Data() map[string]any {
	externalDB := o.Database != "none" && o.Database != "sqlite"
	dialect := sqlDialects[o.Database]
	data := map[string]any{
		"project_name":   o.Name,
		"module_path":    o.module(),
		"go_version":     o.GoVersion,
		"framework":      o.Framework,
		"framework_id":   strings.ReplaceAll(o.Framework, "/", "-"),
		"fiber":          o.Framework == "fiber",
		"database":       o.Database,
		"has_db":         o.Database != "none",
		"external_db":    externalDB,
		"templ":          o.usesTempl(),
		"htmx":           o.Has("htmx"),
		"tailwind":       o.Has("tailwind"),
		"tailwind_cli":   o.Has("tailwind") && !o.Has("react"),
		"react":          o.Has("react"),
		"websocket":      o.Has("websocket"),
		"docker":         o.Has("docker"),
		"github_actions": o.Has("githubaction"),
		"compose":        o.Has("docker") || externalDB,
		"sqlc":           o.Uses("sqlc"),
		"sqlc_engine":    dialect.sqlcEngine,
		"sql_param":      dialect.param,
		"sql_id_type":    dialect.idType,
		// Without goose migrations, sqlc reads the schema from database/schema.sql.
		"sql_schema_file": o.Uses("sqlc") && !o.Uses("goose"),
		"goose":           o.Uses("goose"),
		"goose_driver":    dialect.gooseDriver,
		"gotestsum":       o.Uses("gotestsum"),
		"runner":          o.Runner,
		"runner_make":     o.Runner == "make",
		"runner_just":     o.Runner == "just",
		"runner_task":     o.Runner == "task",
	}
	for k, v := range toolVersions {
		data[k] = v
	}
	return data
}

func oneOf(what, value string, allowed []string) error {
	if slices.Contains(allowed, value) {
		return nil
	}
	return fmt.Errorf("invalid %s %q, allowed: %s", what, value, strings.Join(allowed, ", "))
}
