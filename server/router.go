package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
)

func (server *server) createRouter() {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)

	r.Use(server.renderParamsInit)
	r.Use(server.csrfMiddleware)
	r.Use(server.addCSRFTokenToRenderParams)

	r.Get("/", server.mainGet)
	r.Get("/short", server.mainGet)

	r.Post("/short", server.shortPost)

	r.Get("/i/{id}", server.elongateGet)
	r.Post("/i/{id}", server.elongatePost)

	server.R = r
}

func (server *server) renderParamsInit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mp := make(map[string]interface{})
		ctx := context.WithValue(r.Context(), ctxKeyRenderParms, &mp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *server) addCSRFTokenToRenderParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mp, _ := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})
		(*mp)["csrf_token"] = csrf.Token(r)
		next.ServeHTTP(w, r)
	})
}
