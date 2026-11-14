package server

import (
	"net/http"
	"os"
	"time"
)

type Server struct{}

func NewServer() *http.Server {
	s := &Server{}

	return &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
