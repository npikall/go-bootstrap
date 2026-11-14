//go:build integration

package database

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := mongodb.Run(ctx, "mongo:8")
	if err != nil {
		log.Fatalf("could not start mongo container: %v", err)
	}

	endpoint, err := container.PortEndpoint(ctx, "27017/tcp", "")
	if err != nil {
		log.Fatalf("could not get mongo endpoint: %v", err)
	}
	if host, port, err = net.SplitHostPort(endpoint); err != nil {
		log.Fatalf("could not parse mongo endpoint: %v", err)
	}
	username, password = "", ""

	code := m.Run()
	if err := testcontainers.TerminateContainer(container); err != nil {
		log.Printf("could not terminate mongo container: %v", err)
	}
	os.Exit(code)
}

func TestHealth(t *testing.T) {
	stats := New().Health()

	if stats["status"] != "up" {
		t.Fatalf("status = %q, want %q (error: %s)", stats["status"], "up", stats["error"])
	}
	if stats["message"] != "It's healthy" {
		t.Errorf("message = %q, want %q", stats["message"], "It's healthy")
	}
}

func TestClose(t *testing.T) {
	if err := New().Close(); err != nil {
		t.Fatalf("Close() = %v, want nil", err)
	}
}
