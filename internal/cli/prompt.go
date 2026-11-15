package cli

import (
	"github.com/charmbracelet/huh"

	"github.com/npikall/go-bootstrap/internal/scaffold"
)

// prompt asks for every option that was not given as a flag.
func prompt(opts *scaffold.Options, set map[string]bool) error {
	fields := promptFields(opts, set)
	if len(fields) == 0 {
		return nil
	}
	// One group per field, so the form shows a single question at a time
	// and its title is never cut off.
	groups := make([]*huh.Group, len(fields))
	for i, f := range fields {
		groups[i] = huh.NewGroup(f)
	}
	return huh.NewForm(groups...).Run()
}

// promptFields returns one question per option not given as a flag, keyed by the
// flag name, in the order they are asked.
func promptFields(opts *scaffold.Options, set map[string]bool) []huh.Field {
	questions := []struct {
		flag  string
		field func(key string) huh.Field
	}{
		{"name", func(key string) huh.Field {
			return huh.NewInput().
				Key(key).
				Title("Project name").
				Value(&opts.Name).
				Validate(scaffold.ValidateName)
		}},
		{"module", func(key string) huh.Field {
			return huh.NewInput().
				Key(key).
				Title("Module path").
				Placeholder("defaults to the project name").
				Value(&opts.ModulePath).
				Validate(scaffold.ValidateModulePath)
		}},
		{"framework", func(key string) huh.Field {
			return selectField(key, "Framework", scaffold.Frameworks, &opts.Framework)
		}},
		{"driver", func(key string) huh.Field {
			return selectField(key, "Database driver", scaffold.Databases, &opts.Database)
		}},
		{"frontend", func(key string) huh.Field {
			return huh.NewSelect[string]().
				Key(key).
				Title("Frontend").
				Options(
					huh.NewOption("none (JSON API only)", "none"),
					huh.NewOption("html (html/template from the standard library)", "html"),
					huh.NewOption("templ", "templ"),
					huh.NewOption("react (Vite + TypeScript)", "react"),
				).
				Value(&opts.Frontend)
		}},
		{"feature", func(key string) huh.Field {
			return huh.NewMultiSelect[string]().
				Key(key).
				Title("Advanced features").
				Description("htmx needs the html or templ frontend, tailwind needs a frontend").
				Options(huh.NewOptions(scaffold.Features...)...).
				Value(&opts.Features).
				// The frontend question comes first, so opts.Frontend is final here.
				Validate(func(features []string) error { return scaffold.ValidateFeatures(features, opts.Frontend) })
		}},
		{"tool", func(key string) huh.Field {
			return huh.NewMultiSelect[string]().
				Key(key).
				Title("Development tools").
				Description("sqlc and goose need a SQL database driver").
				Options(huh.NewOptions(scaffold.Tools...)...).
				Value(&opts.Tools).
				// The driver question comes first, so opts.Database is final here.
				Validate(func(tools []string) error { return scaffold.ValidateTools(tools, opts.Database) })
		}},
		{"runner", func(key string) huh.Field {
			return selectField(key, "Task runner", scaffold.Runners, &opts.Runner)
		}},
		{"git", func(key string) huh.Field {
			return selectField(key, "Git", scaffold.GitModes, &opts.Git)
		}},
	}
	var fields []huh.Field
	for _, q := range questions {
		if !set[q.flag] {
			fields = append(fields, q.field(q.flag))
		}
	}
	return fields
}

func selectField(key, title string, options []string, value *string) huh.Field {
	return huh.NewSelect[string]().
		Key(key).
		Title(title).
		Options(huh.NewOptions(options...)...).
		Value(value)
}
