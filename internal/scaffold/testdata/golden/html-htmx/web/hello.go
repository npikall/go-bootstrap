package web

import (
	"embed"
	"html/template"

	"github.com/gofiber/fiber/v2"
)

//go:embed "templates"
var templateFiles embed.FS

// templates holds every page and fragment in web/templates, parsed once at startup.
var templates = template.Must(template.ParseFS(templateFiles, "templates/*.html"))

func HelloFormHandler(c *fiber.Ctx) error {
	return render(c, "base", "")
}

func HelloWebHandler(c *fiber.Ctx) error {
	return render(c, "hello", c.FormValue("name"))
}

func render(c *fiber.Ctx, name string, data any) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return templates.ExecuteTemplate(c.Response().BodyWriter(), name, data)
}
