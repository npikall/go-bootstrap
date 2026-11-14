package database

import (
	"context"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload" // loads .env before the variables below are read
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	Health() map[string]string

	// Close terminates the database connection.
	Close() error
}

type service struct {
	client *mongo.Client
}

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_ROOT_PASSWORD")
)

func New() Service {
	uri := url.URL{Scheme: "mongodb", Host: net.JoinHostPort(host, port)}
	if username != "" {
		uri.User = url.UserPassword(username, password)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri.String()))
	if err != nil {
		log.Fatal(err)
	}
	return &service{client: client}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.client.Ping(ctx, nil); err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}
	return map[string]string{"status": "up", "message": "It's healthy"}
}

func (s *service) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.client.Disconnect(ctx)
}
