package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
)

func (server *Server) createRouter() {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)

	r.Use(server.renderParamsInit)
	r.Use(server.csrfMiddleware)
	r.Use(server.addCSRFTokenToRenderParams)

	r.Get("/", server.mainGet)
	r.Get("/short", server.mainGet)

	r.Post("/short", server.shortPost)

	r.Route("/i/{id}", func(r chi.Router) {
		r.Use(server.urlCtx)
		r.Use(server.parseForm)

		r.Get("/", server.elongateGet)
		r.Post("/", server.elongatePost)
	})

	server.R = r
}

func (server *Server) renderParamsInit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mp := make(map[string]interface{})
		ctx := context.WithValue(r.Context(), ctxKeyRenderParms, &mp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) addCSRFTokenToRenderParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mp, _ := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})
		(*mp)["csrf_token"] = csrf.Token(r)
		next.ServeHTTP(w, r)
	})
}

func (server *Server) parseForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errParseForm := r.ParseForm()
		if errParseForm != nil {
			fmt.Println("parse error", errParseForm)
			_, err := w.Write([]byte(ErrInternalError))
			if err != nil {
				panic(err)
			}
		}

		next.ServeHTTP(w, r)
	})
}
