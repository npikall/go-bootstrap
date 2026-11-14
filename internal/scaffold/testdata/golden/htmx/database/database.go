package database

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // registers the pgx driver

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
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(username, password),
		Host:     net.JoinHostPort(host, port),
		Path:     database,
		RawQuery: url.Values{"sslmode": {"disable"}, "search_path": {schema}}.Encode(),
	}
	db, err := sql.Open("pgx", dsn.String())
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
