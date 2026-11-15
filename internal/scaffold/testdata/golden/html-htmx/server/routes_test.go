package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHelloWorldHandler(t *testing.T) {
	s := &FiberServer{App: fiber.New()}
	s.RegisterFiberRoutes()

	resp, err := s.Test(httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["message"] != "Hello World" {
		t.Errorf("message = %q, want %q", body["message"], "Hello World")
	}
}
