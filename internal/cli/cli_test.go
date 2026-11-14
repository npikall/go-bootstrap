package cli

import (
	"bytes"
	"io"
	"maps"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

func TestParseFlags(t *testing.T) {
	opts, set, err := parseFlags([]string{
		"-n", "demo",
		"--module", "github.com/acme/demo",
		"-f", "chi",
		"--driver", "postgres",
		"--feature", "htmx,tailwind",
		"--feature", "docker",
		"--runner", "just",
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if opts.Name != "demo" || opts.ModulePath != "github.com/acme/demo" || opts.Framework != "chi" ||
		opts.Database != "postgres" || opts.Runner != "just" || opts.Git != "commit" {
		t.Errorf("unexpected options: %+v", opts)
	}
	if want := []string{"htmx", "tailwind", "docker"}; !slices.Equal(opts.Features, want) {
		t.Errorf("features = %v, want %v", opts.Features, want)
	}
	wantSet := []string{"driver", "feature", "framework", "module", "name", "runner"}
	if got := slices.Sorted(maps.Keys(set)); !slices.Equal(got, wantSet) {
		t.Errorf("set flags = %v, want %v", got, wantSet)
	}
}

func TestParseFlagsRejectsInvalidChoices(t *testing.T) {
	_, _, err := parseFlags([]string{"-framework", "rails", "-feature", "react,htmx"}, io.Discard)
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, want := range []string{"invalid framework", "mutually exclusive"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not contain %q", err, want)
		}
	}
}

func TestParseFlagsRejectsArguments(t *testing.T) {
	if _, _, err := parseFlags([]string{"demo"}, io.Discard); err == nil {
		t.Fatal("expected an error for positional arguments")
	}
}

func TestRunVersion(t *testing.T) {
	var stdout bytes.Buffer
	if code := Run([]string{"version"}, fstest.MapFS{}, "v1.2.3", &stdout, io.Discard); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if got := stdout.String(); got != "v1.2.3\n" {
		t.Errorf("stdout = %q, want %q", got, "v1.2.3\n")
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var stderr bytes.Buffer
	if code := Run([]string{"destroy"}, fstest.MapFS{}, "dev", io.Discard, &stderr); code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), `unknown command "destroy"`) {
		t.Errorf("stderr = %q", stderr.String())
	}
}
