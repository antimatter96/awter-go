// Pacakge server stores all the constants used all over the service

package server

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/gorilla/securecookie"

	"github.com/antimatter96/awter-go/db/url"
)

type contextKey int

// ctxKeyRenderParms is key for renderParams in context
const ctxKeyRenderParms contextKey = 1
const ctxKeyURLObject contextKey = 2

type server struct {
	R *chi.Mux

	shortnerTemplate *template.Template
	createdTemplate  *template.Template
	elongateTemplate *template.Template

	cookie *securecookie.SecureCookie

	csrfMiddleware func(http.Handler) http.Handler

	urlService url.Service

	BcryptCost int
}

// Shortner returns a
func Shortner(templatePath string, urlService url.Service) *server {
	shortner := server{urlService: urlService}

	shortner.BcryptCost = 12

	shortner.parseTemplates(templatePath)
	shortner.initCSRF("s6v9y$B&E)H@McQfThWmZq4t7w!z%C*F", true) // Hardcode now
	shortner.createRouter()

	return &shortner
}
