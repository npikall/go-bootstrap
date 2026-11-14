// Package scaffold renders the project template into a new directory.
//
// The template follows copier conventions: path segments are Jinja templates
// and a segment rendering to an empty string skips that file or directory.
// Files ending in ".jinja" are rendered and lose the suffix, all other files
// are copied verbatim. Files below "_includes" are only used via include.
package scaffold

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	minijinja "github.com/mitsuhiko/minijinja/minijinja-go/v2"
	"github.com/mitsuhiko/minijinja/minijinja-go/v2/syntax"
)

const (
	templateSuffix = ".jinja"
	includesDir    = "_includes"
)

// Render writes the rendered template into dest, which must not exist or be empty.
func Render(tmpl fs.FS, opts Options, dest string) error {
	if err := opts.Validate(); err != nil {
		return err
	}
	if err := ensureEmpty(dest); err != nil {
		return err
	}
	files, err := RenderFiles(tmpl, opts.Data())
	if err != nil {
		return err
	}
	for name, content := range files {
		target := filepath.Join(dest, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, content, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// RenderFiles renders the template in memory and returns file contents keyed by output path.
func RenderFiles(tmpl fs.FS, data map[string]any) (map[string][]byte, error) {
	env := newEnvironment(tmpl)
	files := make(map[string][]byte)
	err := fs.WalkDir(tmpl, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if name == includesDir {
			return fs.SkipDir
		}
		if d.IsDir() || name == "." {
			return nil
		}
		out, err := renderPath(env, name, data)
		if err != nil {
			return err
		}
		if out == "" {
			return nil
		}
		content, err := fs.ReadFile(tmpl, name)
		if err != nil {
			return err
		}
		if strings.HasSuffix(name, templateSuffix) {
			content, err = renderContent(env, name, data)
			if err != nil {
				return err
			}
		}
		files[out] = content
		return nil
	})
	return files, err
}

func newEnvironment(tmpl fs.FS) *minijinja.Environment {
	env := minijinja.NewEnvironment()
	env.SetAutoEscapeFunc(func(string) minijinja.AutoEscape { return minijinja.AutoEscapeNone })
	env.SetUndefinedBehavior(minijinja.UndefinedStrict)
	env.SetWhitespace(syntax.WhitespaceConfig{KeepTrailingNewline: true, TrimBlocks: true, LstripBlocks: true})
	env.SetLoader(func(name string) (string, error) {
		b, err := fs.ReadFile(tmpl, name)
		return string(b), err
	})
	return env
}

// renderPath returns the output path for a template file, or "" if it is skipped.
func renderPath(env *minijinja.Environment, name string, data map[string]any) (string, error) {
	segments := strings.Split(strings.TrimSuffix(name, templateSuffix), "/")
	for i, segment := range segments {
		if !strings.ContainsAny(segment, "{%") {
			continue
		}
		t, err := env.TemplateFromString(segment)
		if err != nil {
			return "", fmt.Errorf("parse path %s: %w", name, err)
		}
		rendered, err := t.Render(data)
		if err != nil {
			return "", fmt.Errorf("render path %s: %w", name, err)
		}
		if rendered == "" {
			return "", nil
		}
		segments[i] = rendered
	}
	return path.Join(segments...), nil
}

func renderContent(env *minijinja.Environment, name string, data map[string]any) ([]byte, error) {
	t, err := env.GetTemplate(name)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", name, err)
	}
	out, err := t.Render(data)
	if err != nil {
		return nil, fmt.Errorf("render %s: %w", name, err)
	}
	return []byte(out), nil
}

func ensureEmpty(dir string) error {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("directory %s already exists and is not empty", dir)
	}
	return nil
}
