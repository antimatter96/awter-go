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
	handlers.Init()
}

func main() {

	handlers.Init()
	mainRouter := mux.NewRouter()
	mainRouter.Handle("/", handlers.ShortnerRouter())

	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./template/static/"))))

	output, _ := constants.Value("output").(string)

	file, err := os.Create(output)
	if err != nil {
		fmt.Printf("could not create file %s : %v", output, err)
	}

	loggedRouter := gorillaHandlers.LoggingHandler(file, mainRouter)

	port, _ := constants.Value("port").(string)

	if err := http.ListenAndServe(port, loggedRouter); err != nil {
		fmt.Printf("error starting server: %v", err)
	}
}
