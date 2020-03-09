package server

import (
	"context"
	"fmt"
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

	r.Get("/", mainGet)
	r.Get("/short", server.mainGet)

	r.Post("/short", server.shortPost)

	r.Get("/i/{id}", server.elongateGet)
	r.Post("/i/{id}", server.elongatePost)

	server.R = r
}

func (server *server) mainGet(w http.ResponseWriter, r *http.Request) {
	renderParams := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})
	fmt.Println(renderParams)

	server.shortnerTemplate.Execute(w, renderParams)
}

func (server *server) shortPost(w http.ResponseWriter, r *http.Request) {
	server.shortnerTemplate.Execute(w, nil)
}

func (server *server) elongateGet(w http.ResponseWriter, r *http.Request) {
	server.shortnerTemplate.Execute(w, nil)
}

func (server *server) elongatePost(w http.ResponseWriter, r *http.Request) {
	server.shortnerTemplate.Execute(w, nil)
}

func (server *server) renderParamsInit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mp := make(map[string]interface{})
		ctx := context.WithValue(r.Context(), ctxKeyRenderParms, &mp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func addCSRFTokenToRenderParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mp, _ := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})
		fmt.Println("hello", mp)
		(*mp)["csrf_token"] = csrf.Token(r)
		fmt.Println("hello", mp)
		next.ServeHTTP(w, r)
	})
}
