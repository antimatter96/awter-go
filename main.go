package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/antimatter96/awter-go/cache"
	"github.com/antimatter96/awter-go/constants"
	"github.com/antimatter96/awter-go/db"
	"github.com/antimatter96/awter-go/handlers"

	gorillaHandlers "github.com/gorilla/handlers"
)

func init() {
	config := flag.String("config", "config", "config file")
	store := flag.String("store", "redis", "The store to use:\n\tMySQL(mysql) or\n\tRedis(redis)\n")
	flag.Parse()

	if err := constants.Init(*config); err != nil {
		fmt.Printf("cant initialize constants : %v", err)
	}

	cache.Init(*store)
	db.InitRedis()
	db.InitMySQL()
	handlers.Init(*store)
}

func main() {

	mainRouter := mux.NewRouter().StrictSlash(false)

	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./template/static/"))))

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
