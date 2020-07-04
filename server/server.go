// Pacakge server stores all the constants used all over the service

package server

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/securecookie"

	"github.com/antimatter96/awter-go/customcrypto"
	"github.com/antimatter96/awter-go/db/url"
)

type contextKey int

// ctxKeyRenderParms is key for renderParams in context
const ctxKeyRenderParms contextKey = 1

// ctxKeyRenderParms is key for URLObject in context
const ctxKeyURLObject contextKey = 2

// Server incapsulates all the things
// templates + routes + everything else
type Server struct {
	R *chi.Mux

	shortnerTemplate *template.Template
	createdTemplate  *template.Template
	elongateTemplate *template.Template

	cookie *securecookie.SecureCookie

	csrfMiddleware func(http.Handler) http.Handler

	urlService url.Service

	customcrypto    customcrypto.CustomCrypto
	passwordChecker customcrypto.PasswordChecker
}

// Shortner returns a
func Shortner(templatePath string, urlService url.Service, customcryptoImplementation customcrypto.CustomCrypto, passwordChecker customcrypto.PasswordChecker) *Server {
	shortner := Server{urlService: urlService, customcrypto: customcryptoImplementation, passwordChecker: passwordChecker}

	shortner.parseTemplates(templatePath)
	shortner.initCSRF("s6v9y$B&E)H@McQfThWmZq4t7w!z%C*F", true) // Hardcode now
	shortner.createRouter()

	return &shortner
}
