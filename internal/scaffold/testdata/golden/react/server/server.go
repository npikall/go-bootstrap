package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/acme/demo/database"
)

type FiberServer struct {
	*fiber.App
	db database.Service
}

func New() *FiberServer {
	return &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "demo",
			AppName:      "demo",
			// main logs the listen address, so the banner would only repeat it.
			DisableStartupMessage: true,
		}),
		db: database.New(),
	}
}
