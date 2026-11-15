package cli

import (
	"bytes"
	"io"
	"maps"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/npikall/go-bootstrap/internal/scaffold"
)

func TestParseFlags(t *testing.T) {
	opts, set, err := parseFlags([]string{
		"-n", "demo",
		"--module", "github.com/acme/demo",
		"-f", "chi",
		"--driver", "postgres",
		"--frontend", "html",
		"--feature", "htmx,tailwind",
		"--feature", "docker",
		"--runner", "just",
		"--tool", "sqlc,goose",
		"--tool", "gotestsum",
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if opts.Name != "demo" || opts.ModulePath != "github.com/acme/demo" || opts.Framework != "chi" ||
		opts.Database != "postgres" || opts.Frontend != "html" || opts.Runner != "just" || opts.Git != "commit" {
		t.Errorf("unexpected options: %+v", opts)
	}
	if want := []string{"htmx", "tailwind", "docker"}; !slices.Equal(opts.Features, want) {
		t.Errorf("features = %v, want %v", opts.Features, want)
	}
	if want := []string{"sqlc", "goose", "gotestsum"}; !slices.Equal(opts.Tools, want) {
		t.Errorf("tools = %v, want %v", opts.Tools, want)
	}
	wantSet := []string{"driver", "feature", "framework", "frontend", "module", "name", "runner", "tool"}
	if got := slices.Sorted(maps.Keys(set)); !slices.Equal(got, wantSet) {
		t.Errorf("set flags = %v, want %v", got, wantSet)
	}
}

func TestParseFlagsRejectsInvalidChoices(t *testing.T) {
	_, _, err := parseFlags([]string{"-framework", "rails", "-frontend", "react", "-feature", "htmx", "-tool", "sqlc"}, io.Discard)
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, want := range []string{"invalid framework", "htmx requires", "sqlc requires"} {
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

func TestPromptFieldsAskOnlyForOptionsNotGivenAsFlags(t *testing.T) {
	all := []string{"name", "module", "framework", "driver", "frontend", "feature", "tool", "runner", "git"}
	tests := map[string]struct {
		set  []string
		want []string
	}{
		"no flags":      {set: nil, want: all},
		"all flags":     {set: all, want: nil},
		"only the rest": {set: []string{"name", "driver", "runner"}, want: []string{"module", "framework", "frontend", "feature", "tool", "git"}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			set := make(map[string]bool)
			for _, f := range tt.set {
				set[f] = true
			}
			opts := scaffold.Defaults()
			var got []string
			for _, f := range promptFields(&opts, set) {
				got = append(got, f.GetKey())
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("asked for %v, want %v", got, tt.want)
			}
		})
	}
}
