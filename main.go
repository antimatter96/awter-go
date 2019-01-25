package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"./cache"
	"./constants"
	"./db"
	"./handlers"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
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

	router := httprouter.New()

	router.GET("/short", handlers.Wrapper(handlers.ExtractSessionID(handlers.ShortnerGet)))
	router.POST("/short", handlers.Wrapper(handlers.ExtractSessionID(handlers.ShortnerPost)))

	router.POST("/i/:id", handlers.Wrapper(handlers.ExtractSessionID(handlers.ElongatePost)))
	router.GET("/i/:id", handlers.Wrapper(handlers.ExtractSessionID(handlers.ElongateGet)))

	router.ServeFiles("/static/*filepath", http.Dir("./template/static/"))

	output, _ := constants.Value("output").(string)

	file, err := os.Create(output)
	if err != nil {
		fmt.Printf("could not create file %s : %v", output, err)
	}

	loggedRouter := gorillaHandlers.LoggingHandler(file, router)

	port, _ := constants.Value("port").(string)

	if err := http.ListenAndServe(port, loggedRouter); err != nil {
		fmt.Printf("error starting server: %v", err)
	}
}
