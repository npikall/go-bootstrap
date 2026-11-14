package server

import (
	"net/http"
	"os"
	"time"

	"github.com/acme/demo/database"
)

type Server struct {
	db database.Service
}

func NewServer() *http.Server {
	s := &Server{db: database.New()}

	return &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
