package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	"github.com/acme/demo/web"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.Use(cors.New())

	s.Get("/", s.HelloWorldHandler)

	s.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.Files),
		PathPrefix: "assets",
	}))
	s.Get("/web", web.HelloFormHandler)
	s.Post("/hello", web.HelloWebHandler)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello World"})
}
