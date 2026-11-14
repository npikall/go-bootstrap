package server

import (
	"fmt"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func (s *FiberServer) websocketHandler(conn *websocket.Conn) {
	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		if err := conn.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
}
