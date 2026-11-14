package cli

import (
	"github.com/charmbracelet/huh"

	"github.com/npikall/go-bootstrap/internal/scaffold"
)

// prompt asks for every option that was not given as a flag.
func prompt(opts *scaffold.Options, set map[string]bool) error {
	var fields []huh.Field
	if !set["name"] {
		fields = append(fields, huh.NewInput().
			Title("Project name").
			Value(&opts.Name).
			Validate(scaffold.ValidateName))
	}
	if !set["module"] {
		fields = append(fields, huh.NewInput().
			Title("Module path").
			Placeholder("defaults to the project name").
			Value(&opts.ModulePath).
			Validate(scaffold.ValidateModulePath))
	}
	if !set["framework"] {
		fields = append(fields, selectField("Framework", scaffold.Frameworks, &opts.Framework))
	}
	if !set["driver"] {
		fields = append(fields, selectField("Database driver", scaffold.Databases, &opts.Database))
	}
	if !set["feature"] {
		fields = append(fields, huh.NewMultiSelect[string]().
			Title("Advanced features").
			Options(huh.NewOptions(scaffold.Features...)...).
			Value(&opts.Features).
			Validate(scaffold.ValidateFeatures))
	}
	if !set["runner"] {
		fields = append(fields, selectField("Task runner", scaffold.Runners, &opts.Runner))
	}
	if !set["git"] {
		fields = append(fields, selectField("Git", scaffold.GitModes, &opts.Git))
	}
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

func selectField(title string, options []string, value *string) huh.Field {
	return huh.NewSelect[string]().
		Title(title).
		Options(huh.NewOptions(options...)...).
		Value(value)
}
