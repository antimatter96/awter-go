package handlers

import (
	"./common"
	"./shortner"
	"github.com/gorilla/mux"
)

func Init() {
	common.InitCommon()
	shortner.InitShortner()
}

func ShortnerRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/short", shortner.ShortnerGet)
	r.HandleFunc("/short", shortner.ShortnerPost)
	r.HandleFunc("/i/:id", shortner.ElongatePost)
	r.HandleFunc("/i/:id", shortner.ElongateGet)
	return r
}
