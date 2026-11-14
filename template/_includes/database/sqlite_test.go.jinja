package database

import (
	"path/filepath"
	"testing"
)

func newTestService(t *testing.T) Service {
	t.Helper()
	dburl = filepath.Join(t.TempDir(), "test.db")
	return New()
}

func TestHealth(t *testing.T) {
	srv := newTestService(t)
	defer srv.Close()

	stats := srv.Health()
	if stats["status"] != "up" {
		t.Fatalf("status = %q, want %q (error: %s)", stats["status"], "up", stats["error"])
	}
}

func TestClose(t *testing.T) {
	if err := newTestService(t).Close(); err != nil {
		t.Fatalf("Close() = %v, want nil", err)
	}
}
