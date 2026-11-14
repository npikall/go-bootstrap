package server

import (
	"net/http"
	"time"

	"github.com/acme/demo/database"
)

type Server struct {
	db database.Service
}

func NewServer() *http.Server {
	s := &Server{db: database.New()}

	return &http.Server{
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
