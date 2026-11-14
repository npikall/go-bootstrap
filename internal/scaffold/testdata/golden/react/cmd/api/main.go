package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/acme/demo/server"
)

func gracefulShutdown(srv *server.FiberServer, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}
	log.Println("Server exiting")
	done <- true
}

func main() {
	ln, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatalln("listen error:", err)
	}
	log.Printf("listening on %s", ln.Addr())

	srv := server.New()
	srv.RegisterFiberRoutes()

	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)

	if err := srv.Listener(ln); err != nil {
		log.Fatalln("http server error:", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
