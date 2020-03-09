package main

import (
	"flag"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/antimatter96/awter-go/server"
)

func main() {
	var port = flag.Int("port", 8080, "port")
	var templatePath = flag.String("template", "./template", "the template directory")

	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	shortner := server.Shortner(*templatePath)

	r.Mount("/", shortner.R)

	http.ListenAndServe(":"+strconv.Itoa(*port), r)
}
