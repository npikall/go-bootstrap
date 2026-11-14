package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload" // loads .env before the variables below are read
	_ "modernc.org/sqlite"                // registers the sqlite driver

	"github.com/acme/demo/database/store"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	Health() map[string]string

	// Close terminates the database connection.
	Close() error

	// Queries returns the queries sqlc generates from database/queries.
	Queries() *store.Queries
}

type service struct {
	db *sql.DB
}

var (
	dburl      = os.Getenv("DB_URL")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite", dburl)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{db: db}
	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}

	stats := s.db.Stats()
	return map[string]string{
		"status":           "up",
		"message":          "It's healthy",
		"open_connections": strconv.Itoa(stats.OpenConnections),
		"in_use":           strconv.Itoa(stats.InUse),
		"idle":             strconv.Itoa(stats.Idle),
	}
}

func (s *service) Close() error {
	dbInstance = nil
	return s.db.Close()
}

func (s *service) Queries() *store.Queries {
	return store.New(s.db)
}
