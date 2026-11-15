package server

import (
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
}

func New() *FiberServer {
	return &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "demo",
			AppName:      "demo",
			// main logs the listen address, so the banner would only repeat it.
			DisableStartupMessage: true,
		}),
	}
}
