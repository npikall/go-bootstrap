package server

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.Use(cors.New())

	s.Get("/", s.HelloWorldHandler)
	s.Get("/health", s.healthHandler)
	s.Get("/websocket", websocket.New(s.websocketHandler))
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello World"})
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
