package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("could not open websocket: %v", err)
		return
	}
	defer func() { _ = conn.CloseNow() }()

	ctx := conn.CloseRead(r.Context())
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		if err := conn.Write(ctx, websocket.MessageText, []byte(payload)); err != nil {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
