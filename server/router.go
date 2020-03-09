package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (server *server) createRouter() {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)

	r.Get("/", server.mainGet)
	r.Get("/short", server.mainGet)

	r.Post("/short", server.shortPost)

	r.Get("/i/{id}", server.elongateGet)
	r.Post("/i/{id}", server.elongatePost)

	server.R = r
}

func (server *server) mainGet(w http.ResponseWriter, r *http.Request) {
	server.shortnerTemplate.Execute(w, nil)
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
