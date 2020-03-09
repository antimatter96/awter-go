package server

import (
	"text/template"

	"github.com/go-chi/chi"
	"github.com/gorilla/securecookie"
)

type server struct {
	R *chi.Mux

	shortnerTemplate *template.Template
	createdTemplate  *template.Template
	elongateTemplate *template.Template

	Cookie *securecookie.SecureCookie
}

// Shortner returns a
func Shortner(templatePath string) *server {
	shortner := server{}

	shortner.parseTemplates(templatePath)
	shortner.createRouter()

	return &shortner
}
