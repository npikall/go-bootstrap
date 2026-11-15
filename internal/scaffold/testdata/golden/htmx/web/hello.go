package web

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
)

var helloForm = templ.Handler(HelloForm())

func HelloFormHandler(w http.ResponseWriter, r *http.Request) {
	helloForm.ServeHTTP(w, r)
}

func HelloWebHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := HelloPost(r.FormValue("name")).Render(r.Context(), w); err != nil {
		log.Printf("Error rendering in HelloWebHandler: %v", err)
	}
}
