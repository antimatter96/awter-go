package handlers

import (
	"fmt"
	"net/http"

	"github.com/antimatter96/awter-go/handlers/common"
	"github.com/antimatter96/awter-go/handlers/shortner"
	"github.com/gorilla/mux"
)

func Init() {
	common.InitCommon()
	shortner.InitShortner()
}
func notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "FUCK")
}

func ShortnerRouter(r *mux.Router) {
	r.HandleFunc("/", shortner.Get).Methods("GET")
	r.HandleFunc("/short", shortner.Get).Methods("GET")
	r.HandleFunc("/short", shortner.Post).Methods("POST")
	r.HandleFunc("/i/{id}", shortner.ElongatePost).Methods("GET")
	r.HandleFunc("/i/{id}", shortner.ElongateGet).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(notFound)
}
