// Package cli implements the go-bootstrap command line interface.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/npikall/go-bootstrap/internal/scaffold"
)

const usage = `go-bootstrap scaffolds new Go projects.

Usage:
  go-bootstrap create [flags]
  go-bootstrap version

Flags of create:
`

// Run executes the CLI and returns the process exit code.
func Run(args []string, tmpl fs.FS, version string, stdout, stderr io.Writer) int {
	if err := run(args, tmpl, version, stdout, stderr); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			_, _ = fmt.Fprintln(stderr, "error:", err)
		}
		return 1
	}
	return 0
}

func run(args []string, tmpl fs.FS, version string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		_, _ = fmt.Fprint(stderr, usage)
		newFlagSet(&scaffold.Options{}, stderr).PrintDefaults()
		return flag.ErrHelp
	}
	switch args[0] {
	case "version":
		_, err := fmt.Fprintln(stdout, version)
		return err
	case "create":
		return create(args[1:], tmpl, stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func create(args []string, tmpl fs.FS, stdout, stderr io.Writer) error {
	opts, set, err := parseFlags(args, stderr)
	if err != nil {
		return err
	}
	if isInteractive() {
		if err := prompt(&opts, set); err != nil {
			return err
		}
	} else if !set["name"] {
		return errors.New("flag -name is required when not running in a terminal")
	}
	opts.GoVersion = scaffold.ToolchainGoVersion()
	if err := scaffold.Render(tmpl, opts, opts.Name); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(stdout, "Created %s, finishing setup:\n", opts.Name)
	if err := scaffold.Finalize(opts.Name, opts, stdout); err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "\nDone. Next steps:\n  cd %s\n  %s run\n", opts.Name, opts.Runner)
	return err
}

// parseFlags returns the options and the set of canonical flag names given explicitly.
func parseFlags(args []string, output io.Writer) (scaffold.Options, map[string]bool, error) {
	opts := scaffold.Defaults()
	fset := newFlagSet(&opts, output)
	if err := fset.Parse(args); err != nil {
		return opts, nil, err
	}
	if fset.NArg() > 0 {
		return opts, nil, fmt.Errorf("unexpected arguments: %s", strings.Join(fset.Args(), " "))
	}
	aliases := map[string]string{"n": "name", "f": "framework", "d": "driver", "g": "git"}
	set := make(map[string]bool)
	fset.Visit(func(f *flag.Flag) {
		name := f.Name
		if canonical, ok := aliases[name]; ok {
			name = canonical
		}
		set[name] = true
	})
	return opts, set, opts.ValidateChoices()
}

func newFlagSet(opts *scaffold.Options, output io.Writer) *flag.FlagSet {
	fset := flag.NewFlagSet("create", flag.ContinueOnError)
	fset.SetOutput(output)
	stringFlag := func(p *string, name, short, help string) {
		fset.StringVar(p, name, *p, help)
		if short != "" {
			fset.StringVar(p, short, *p, "shorthand for -"+name)
		}
	}
	stringFlag(&opts.Name, "name", "n", "project name, also the output directory")
	stringFlag(&opts.ModulePath, "module", "", "go module path (default: project name)")
	stringFlag(&opts.Framework, "framework", "f", "framework: "+strings.Join(scaffold.Frameworks, ", "))
	stringFlag(&opts.Database, "driver", "d", "database driver: "+strings.Join(scaffold.Databases, ", "))
	stringFlag(&opts.Runner, "runner", "", "task runner: "+strings.Join(scaffold.Runners, ", "))
	stringFlag(&opts.Git, "git", "g", "git setup: "+strings.Join(scaffold.GitModes, ", "))
	listFlag := func(p *[]string, name, help string, allowed []string) {
		fset.Func(name, help+", repeatable or comma separated: "+strings.Join(allowed, ", "), func(v string) error {
			for item := range strings.SplitSeq(v, ",") {
				if item = strings.TrimSpace(item); item != "" {
					*p = append(*p, item)
				}
			}
			return nil
		})
	}
	listFlag(&opts.Features, "feature", "advanced feature", scaffold.Features)
	listFlag(&opts.Tools, "tool", "development tool", scaffold.Tools)
	return fset
}

func isInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}
