//go:build integration

package database

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	endpoint, err := container.PortEndpoint(ctx, "5432/tcp", "")
	if err != nil {
		log.Fatalf("could not get postgres endpoint: %v", err)
	}
	if host, port, err = net.SplitHostPort(endpoint); err != nil {
		log.Fatalf("could not parse postgres endpoint: %v", err)
	}
	database, username, password, schema = "test", "test", "test", "public"

	code := m.Run()
	if err := testcontainers.TerminateContainer(container); err != nil {
		log.Printf("could not terminate postgres container: %v", err)
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
