package web

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed "templates"
var templateFiles embed.FS

// templates holds every page and fragment in web/templates, parsed once at startup.
var templates = template.Must(template.ParseFS(templateFiles, "templates/*.html"))

func HelloFormHandler(w http.ResponseWriter, _ *http.Request) {
	render(w, "base", "")
}

func HelloWebHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	render(w, "base", r.FormValue("name"))
}

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("Error rendering template %s: %v", name, err)
	}
}
