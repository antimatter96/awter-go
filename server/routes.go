package server

import "net/http"

func (server *server) mainGet(w http.ResponseWriter, r *http.Request) {
	renderParams := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})
	server.shortnerTemplate.Execute(w, renderParams)
}

func (server *server) shortPost(w http.ResponseWriter, r *http.Request) {
	server.shortnerTemplate.Execute(w, nil)
}

func (server *server) elongateGet(w http.ResponseWriter, r *http.Request) {
	server.elongateTemplate.Execute(w, nil)
}

func (server *server) elongatePost(w http.ResponseWriter, r *http.Request) {
	server.shortnerTemplate.Execute(w, nil)
}
