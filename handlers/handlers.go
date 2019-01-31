package handlers

import (
	"fmt"
	"net/http"

	"github.com/antimatter96/awter-go/constants"
	"github.com/antimatter96/awter-go/handlers/common"
	"github.com/antimatter96/awter-go/handlers/shortner"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

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
	r.HandleFunc("/", shortner.Get).Methods("GET")
	r.HandleFunc("/short", shortner.Get).Methods("GET")
	r.HandleFunc("/short", shortner.Post).Methods("POST")
	r.HandleFunc("/i/{id}", shortner.ElongatePost).Methods("GET")
	r.HandleFunc("/i/{id}", shortner.ElongateGet).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(notFound)
}
