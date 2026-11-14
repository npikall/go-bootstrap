package scaffold

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
	"golang.org/x/mod/module"
)

var update = flag.Bool("update", false, "update golden files")

var templateDir = filepath.Join("..", "..", "template")

func testOptions(mutate func(*Options)) Options {
	opts := Defaults()
	opts.Name = "demo"
	opts.ModulePath = "github.com/acme/demo"
	opts.GoVersion = "1.27.4"
	if mutate != nil {
		mutate(&opts)
	}
	return opts
}

var featureSets = map[string][]string{
	"none":     nil,
	"htmx":     {"htmx", "tailwind", "websocket", "docker", "githubaction"},
	"tailwind": {"tailwind"},
	"react":    {"react", "tailwind", "docker", "githubaction"},
}

// matrix covers every framework, database, runner and the representative feature sets.
func matrix() map[string]Options {
	cases := make(map[string]Options)
	for _, fw := range Frameworks {
		for _, db := range Databases {
			for setName, features := range featureSets {
				for _, runner := range Runners {
					name := fmt.Sprintf("%s/%s/%s/%s", fw, db, setName, runner)
					cases[name] = testOptions(func(o *Options) {
						o.Framework, o.Database, o.Features, o.Runner = fw, db, features, runner
					})
				}
			}
		}
	}
	return cases
}

func TestRenderMatrix(t *testing.T) {
	tmpl := os.DirFS(templateDir)
	for name, opts := range matrix() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			files, err := RenderFiles(tmpl, opts.Data())
			if err != nil {
				t.Fatal(err)
			}
			checkFileSet(t, opts, files)
			for p, content := range files {
				checkContent(t, p, content)
			}
		})
	}
}

func checkFileSet(t *testing.T, opts Options, files map[string][]byte) {
	t.Helper()
	data := opts.Data()
	has := func(prefix string) bool {
		for p := range files {
			if p == prefix || strings.HasPrefix(p, prefix+"/") {
				return true
			}
		}
		return false
	}
	expect := map[string]bool{
		"go.mod":               true,
		"cmd/api/main.go":      true,
		"server/routes.go":     true,
		".golangci.yml":        true,
		"Makefile":             opts.Runner == "make",
		"justfile":             opts.Runner == "just",
		"Taskfile.yml":         opts.Runner == "task",
		"database":             data["has_db"].(bool),
		"web":                  data["templ"].(bool),
		"web/styles/input.css": data["tailwind_cli"].(bool),
		"frontend":             opts.Has("react"),
		"Dockerfile":           opts.Has("docker"),
		"compose.yaml":         data["compose"].(bool),
		".github":              opts.Has("githubaction"),
		".goreleaser.yml":      opts.Has("githubaction"),
		"server/websocket.go":  opts.Has("websocket"),
	}
	for p, want := range expect {
		if got := has(p); got != want {
			t.Errorf("%s present = %v, want %v", p, got, want)
		}
	}
	for p := range files {
		if strings.HasPrefix(p, includesDir) || strings.HasSuffix(p, templateSuffix) {
			t.Errorf("template internals leaked into output: %s", p)
		}
	}
}

// checkContent enforces formatting rules on every rendered text file.
func checkContent(t *testing.T, name string, content []byte) {
	t.Helper()
	if strings.HasSuffix(name, ".min.js") {
		return
	}
	text := string(content)
	switch {
	case text == "":
		t.Errorf("%s: empty file", name)
		return
	case strings.HasPrefix(text, "\n"):
		t.Errorf("%s: starts with a blank line", name)
	case !strings.HasSuffix(text, "\n") || strings.HasSuffix(text, "\n\n"):
		t.Errorf("%s: must end with exactly one newline", name)
	}
	for i, line := range strings.Split(text, "\n") {
		if strings.TrimRight(line, " \t") != line {
			t.Errorf("%s:%d: trailing whitespace", name, i+1)
		}
		if strings.Contains(line, "{%") || strings.Contains(line, "%}") || strings.Contains(line, "{#") {
			t.Errorf("%s:%d: unrendered jinja syntax: %s", name, i+1, line)
		}
		if path.Base(name) == "Makefile" && strings.HasPrefix(line, " ") {
			t.Errorf("%s:%d: Makefile lines must not be indented with spaces", name, i+1)
		}
	}
	if before, _, found := strings.Cut(text, "\n\n\n"); found {
		t.Errorf("%s:%d: consecutive blank lines", name, strings.Count(before, "\n")+2)
	}

	switch ext := path.Ext(name); ext {
	case ".go":
		formatted, err := format.Source(content)
		if err != nil {
			t.Errorf("%s: invalid go: %v", name, err)
		} else if !bytes.Equal(formatted, content) {
			t.Errorf("%s: not gofmt formatted, diff:\n%s", name, lineDiff(text, string(formatted)))
		}
	case ".yml", ".yaml":
		var v any
		if err := yaml.Unmarshal(content, &v); err != nil {
			t.Errorf("%s: invalid yaml: %v", name, err)
		}
	case ".json":
		if !json.Valid(content) {
			t.Errorf("%s: invalid json", name)
		}
	}
}

func lineDiff(got, want string) string {
	g, w := strings.Split(got, "\n"), strings.Split(want, "\n")
	for i := range min(len(g), len(w)) {
		if g[i] != w[i] {
			return fmt.Sprintf("line %d\n got: %q\nwant: %q", i+1, g[i], w[i])
		}
	}
	return fmt.Sprintf("line counts differ: got %d, want %d", len(g), len(w))
}

var goldenCases = map[string]Options{
	"minimal": testOptions(nil),
	"htmx": testOptions(func(o *Options) {
		o.Framework, o.Database, o.Runner = "chi", "postgres", "just"
		o.Features = featureSets["htmx"]
	}),
	"react": testOptions(func(o *Options) {
		o.Framework, o.Database, o.Runner = "fiber", "mongo", "task"
		o.Features = append(slices.Clone(featureSets["react"]), "websocket")
	}),
}

// TestGolden compares complete rendered projects with testdata/golden.
// Run "go test ./internal/scaffold -run TestGolden -update" after changing the template.
func TestGolden(t *testing.T) {
	tmpl := os.DirFS(templateDir)
	for name, opts := range goldenCases {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join("testdata", "golden", name)
			files, err := RenderFiles(tmpl, opts.Data())
			if err != nil {
				t.Fatal(err)
			}
			delete(files, "web/assets/js/htmx.min.js")
			// The golden project's own .gitignore ignores .env, so store it under a name that gets committed.
			files[".envrc"] = files[".env"]
			delete(files, ".env")
			if *update {
				writeGolden(t, dir, files)
			}
			want := readGolden(t, dir)
			for p, content := range files {
				if w, ok := want[p]; !ok {
					t.Errorf("unexpected file %s", p)
				} else if !bytes.Equal(content, w) {
					t.Errorf("%s differs from golden file:\n%s", p, lineDiff(string(content), string(w)))
				}
			}
			for p := range want {
				if _, ok := files[p]; !ok {
					t.Errorf("missing file %s", p)
				}
			}
		})
	}
}

func writeGolden(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	for p, content := range files {
		target := filepath.Join(dir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func readGolden(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	golden := os.DirFS(dir)
	files := make(map[string][]byte)
	err := fs.WalkDir(golden, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		files[p], err = fs.ReadFile(golden, p)
		return err
	})
	if err != nil {
		t.Fatalf("read golden files (run with -update to create them): %v", err)
	}
	return files
}

// Template paths end up in the module zip, so they must be valid there for "go install" to work.
func TestTemplatePathsAreModuleSafe(t *testing.T) {
	err := filepath.WalkDir(templateDir, func(p string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(filepath.Join(templateDir, ".."), p)
		if err != nil {
			return err
		}
		if err := module.CheckFilePath(filepath.ToSlash(rel)); err != nil {
			t.Error(err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRenderWritesFiles(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "demo")
	if err := Render(os.DirFS(templateDir), testOptions(nil), dest); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(dest, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "module github.com/acme/demo\n\ngo 1.27.4\n"; string(content) != want {
		t.Errorf("go.mod = %q, want %q", content, want)
	}
}

func TestRenderRefusesNonEmptyDirectory(t *testing.T) {
	dest := t.TempDir()
	if err := os.WriteFile(filepath.Join(dest, "keep"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Render(os.DirFS(templateDir), testOptions(nil), dest); err == nil {
		t.Fatal("expected an error for a non-empty directory")
	}
}

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		mutate  func(*Options)
		wantErr string
	}{
		"valid":             {mutate: nil},
		"empty name":        {mutate: func(o *Options) { o.Name = "" }, wantErr: "invalid project name"},
		"bad module":        {mutate: func(o *Options) { o.ModulePath = "Not A Module" }, wantErr: "invalid module path"},
		"unknown framework": {mutate: func(o *Options) { o.Framework = "rails" }, wantErr: "invalid framework"},
		"unknown driver":    {mutate: func(o *Options) { o.Database = "oracle" }, wantErr: "invalid database driver"},
		"unknown feature":   {mutate: func(o *Options) { o.Features = []string{"grpc"} }, wantErr: "invalid feature"},
		"unknown runner":    {mutate: func(o *Options) { o.Runner = "bazel" }, wantErr: "invalid runner"},
		"react and htmx":    {mutate: func(o *Options) { o.Features = []string{"react", "htmx"} }, wantErr: "mutually exclusive"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := testOptions(tt.mutate).Validate()
			switch {
			case tt.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)):
				t.Fatalf("error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
