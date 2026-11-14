package scaffold

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Finalize resolves dependencies, generates code and initialises git in dir.
func Finalize(dir string, opts Options, log io.Writer) error {
	if err := runSteps(dir, log, []string{"go", "mod", "tidy"}); err != nil {
		return err
	}
	var steps [][]string
	if opts.usesTempl() {
		// Use the templ CLI matching the library version that "go mod tidy" resolved.
		version, err := output(dir, "go", "list", "-m", "-f", "{{.Version}}", "github.com/a-h/templ")
		if err != nil {
			return err
		}
		steps = append(steps, []string{"go", "run", "github.com/a-h/templ/cmd/templ@" + version, "generate"})
	}
	if opts.Git != "skip" {
		steps = append(steps, []string{"git", "init", "--quiet"})
	}
	if opts.Git == "stage" || opts.Git == "commit" {
		steps = append(steps, []string{"git", "add", "--all"})
	}
	if opts.Git == "commit" {
		steps = append(steps, []string{"git", "commit", "--quiet", "--message", "Initial commit"})
	}
	return runSteps(dir, log, steps...)
}

func runSteps(dir string, log io.Writer, steps ...[]string) error {
	for _, step := range steps {
		if _, err := fmt.Fprintf(log, "  %s\n", strings.Join(step, " ")); err != nil {
			return err
		}
		if _, err := output(dir, step[0], step[1:]...); err != nil {
			return err
		}
	}
	return nil
}

func output(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
