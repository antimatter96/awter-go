package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"./cache"
	"./constants"
	"./db"
	"./handlers"

	gorillaHandlers "github.com/gorilla/handlers"
)

func init() {
	config := flag.String("config", "config", "config file")
	flag.Parse()

	if err := constants.Init(*config); err != nil {
		fmt.Printf("cant initialize constants : %v", err)
	}

	cache.Init()
	db.InitRedis()
	db.InitMySQL()
	handlers.Init()
}

func main() {

	mainRouter := mux.NewRouter().StrictSlash(false)
	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./template/static/"))))

	shortnerRouter := mainRouter.PathPrefix("/").Subrouter()
	handlers.ShortnerRouter(shortnerRouter)

	output, _ := constants.Value("output").(string)

	file, err := os.OpenFile(output, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("could not create file %s : %v", output, err)
	}

	loggedRouter := gorillaHandlers.LoggingHandler(file, mainRouter)
	http.Handle("/", mainRouter)

	port, _ := constants.Value("port").(string)

	if err := http.ListenAndServe(port, loggedRouter); err != nil {
		fmt.Printf("error starting server: %v", err)
	}
}
