package server

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/gorilla/securecookie"
)

type contextKey int

// ctxKeyRenderParms is key for renderParams in context
const ctxKeyRenderParms contextKey = 1

type server struct {
	R *chi.Mux

	shortnerTemplate *template.Template
	createdTemplate  *template.Template
	elongateTemplate *template.Template

	cookie *securecookie.SecureCookie

	csrfMiddleware func(http.Handler) http.Handler
}

// Shortner returns a
func Shortner(templatePath string) *server {
	shortner := server{}

	shortner.parseTemplates(templatePath)
	shortner.initCSRF("s6v9y$B&E)H@McQfThWmZq4t7w!z%C*F", true) // Hardcode now
	shortner.createRouter()

	return &shortner
}
