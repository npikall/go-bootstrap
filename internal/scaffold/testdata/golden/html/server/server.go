package server

import (
	"net/http"
	"time"
)

type Server struct{}

func NewServer() *http.Server {
	s := &Server{}

	return &http.Server{
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
