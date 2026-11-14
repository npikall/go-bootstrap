package scaffold

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// TestEndToEnd generates real projects and builds, vets, tests and lints them.
// It needs network access and is therefore opt-in: GO_BOOTSTRAP_E2E=1 go test ./...
// Checks needing golangci-lint, zizmor, just, task or npm are skipped when the tool is missing.
func TestEndToEnd(t *testing.T) {
	if os.Getenv("GO_BOOTSTRAP_E2E") == "" {
		t.Skip("set GO_BOOTSTRAP_E2E=1 to run the end-to-end tests")
	}

	type e2eCase struct {
		opts       Options
		runTargets bool // also run the build and test targets of the task runner
	}
	cases := map[string]e2eCase{}
	add := func(name, fw, db, runner string, features []string, runTargets bool) {
		opts := testOptions(func(o *Options) {
			o.Name, o.Framework, o.Database, o.Runner, o.Features, o.Git = name, fw, db, runner, features, "skip"
			o.GoVersion = ToolchainGoVersion() // the generated go.mod must accept the toolchain that builds it
		})
		cases[name] = e2eCase{opts: opts, runTargets: runTargets}
	}
	for i, fw := range Frameworks {
		add("full-"+fw[:3], fw, Databases[i], Runners[i%len(Runners)], featureSets["htmx"], false)
		add("bare-"+fw[:3], fw, "none", "make", nil, false)
	}
	add("react-fiber", "fiber", "sqlite", "task", append(slices.Clone(featureSets["react"]), "websocket"), true)
	add("react-gin", "gin", "redis", "just", featureSets["react"], false)
	add("tailwind-make", "standard-library", "sqlite", "make", featureSets["tailwind"], true)
	add("tailwind-just", "echo", "none", "just", featureSets["tailwind"], true)
	add("tailwind-task", "httprouter", "none", "task", featureSets["tailwind"], true)

	tmpl := os.DirFS(templateDir)
	root := t.TempDir()
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := tc.opts
			dir := filepath.Join(root, name)
			if err := Render(tmpl, opts, dir); err != nil {
				t.Fatal(err)
			}
			if err := Finalize(dir, opts, io.Discard); err != nil {
				t.Fatal(err)
			}
			mustRun(t, dir, "go", "vet", "-tags", "integration", "./...")
			mustRun(t, dir, "go", "test", "./...")
			if opts.Database == "none" {
				checkFirstLineIsListenAddress(t, dir)
			}
			if _, err := exec.LookPath("golangci-lint"); err == nil {
				mustRun(t, dir, "golangci-lint", "run", "--allow-parallel-runners", "./...")
			}
			if _, err := exec.LookPath("zizmor"); err == nil && opts.Has("githubaction") {
				mustRun(t, dir, "zizmor", "--offline", "--min-severity", "low", ".github")
			}
			if _, err := exec.LookPath(opts.Runner); err == nil && tc.runTargets {
				mustRun(t, dir, opts.Runner, "build")
				mustRun(t, dir, opts.Runner, "test")
			}
			if _, err := exec.LookPath("npm"); err == nil && opts.Has("react") && tc.runTargets {
				mustRun(t, filepath.Join(dir, "frontend"), "npm", "install", "--no-audit", "--no-fund")
				mustRun(t, filepath.Join(dir, "frontend"), "npm", "run", "build")
			}
		})
	}
}

// checkFirstLineIsListenAddress starts the generated server and requires that the
// first thing it prints tells the user where it is listening.
func checkFirstLineIsListenAddress(t *testing.T, dir string) {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "api")
	mustRun(t, dir, "go", "build", "-o", bin, "./cmd/api")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	cmd := exec.Command(bin)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PORT=0")
	cmd.Stdout, cmd.Stderr = w, w
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()

	line := make(chan string, 1)
	go func() {
		s := bufio.NewScanner(r)
		s.Scan()
		line <- s.Text()
	}()
	select {
	case got := <-line:
		if !strings.Contains(got, "listening on ") {
			t.Errorf("first output line = %q, want the listen address", got)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("server printed nothing within 30s")
	}
}

func mustRun(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
}
