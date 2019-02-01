package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/antimatter96/awter-go/constants"
	"github.com/antimatter96/awter-go/handlers/common"
	"github.com/antimatter96/awter-go/handlers/shortner"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func contextInitializer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mp := make(map[string]interface{})
		ctx := context.WithValue(r.Context(), common.CtxKeyResParms, mp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func addCSRFTokenToRenderParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mp, err := r.Context().Value(common.CtxKeyResParms).(map[string]interface{})
		if !err {
			panic("Context is not a map")
		}
		mp["csrf_token"] = csrf.Token(r)
		ctx := context.WithValue(r.Context(), common.CtxKeyResParms, mp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Init(store string) {
	common.InitCommon()
	shortner.InitShortner(store)
}
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "FUCK from Shortner")
}

func ShortnerRouter(r *mux.Router) {
	csrfMiddleware := csrf.Protect(
		[]byte(constants.Value("csrf-token").(string)),
		csrf.FieldName("_csrf_token"),
		csrf.CookieName("_csrf_token"),
		csrf.Secure(constants.ENVIRONMENT != "dev"),
	)

	r.Use(csrfMiddleware)
	r.Use(contextInitializer)
	r.Use(addCSRFTokenToRenderParams)

	r.HandleFunc("/", shortner.Get).Methods("GET")
	r.HandleFunc("/short", shortner.Get).Methods("GET")
	r.HandleFunc("/short", shortner.Post).Methods("POST")
	r.HandleFunc("/i/{id}", shortner.ElongatePost).Methods("GET")
	r.HandleFunc("/i/{id}", shortner.ElongateGet).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(notFound)
}
