package handlers

import (
	"context"
	"fmt"
	"net/http"

	. "github.com/antimatter96/awter-go/handlers/common"
	"github.com/antimatter96/awter-go/handlers/shortner"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func contextInitializer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mp := make(map[string]interface{})
		ctx := context.WithValue(r.Context(), CtxKeyRenderParms, mp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "FUCK from Shortner")
}

var csrfErrorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s", "CSRF SHIT HAPPENED", csrf.FailureReason(r))
})

func Init(store string) {
	InitCommon()
	InitCSRF(csrfErrorHandler)
	shortner.InitShortner(store)
}

func ShortnerRouter(r *mux.Router) {

	r.Use(CSRFMiddleware)
	r.Use(contextInitializer)
	r.Use(addCSRFTokenToRenderParams)

	r.HandleFunc("/", shortner.Get).Methods("GET")
	r.HandleFunc("/short", shortner.Get).Methods("GET")
	r.HandleFunc("/short", shortner.Post).Methods("POST")
	r.HandleFunc("/i/{id}", shortner.ElongatePost).Methods("GET")
	r.HandleFunc("/i/{id}", shortner.ElongateGet).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(notFound)
}
