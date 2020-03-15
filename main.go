package main

import (
	"flag"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/antimatter96/awter-go/db"
	"github.com/antimatter96/awter-go/db/url"
	"github.com/antimatter96/awter-go/server"
)

func main() {
	var port = flag.Int("port", 8080, "port")
	var templatePath = flag.String("template", "./template", "the template directory")
	var mySQLConnectionString = flag.String("mysqlURL", "user:password@/name?parseTime=true", "MySQL connection string")
	var redisAddressstring = flag.String("redisURL", "127.0.0.1:6379", "redis connection string")

	flag.Parse()

	var urlSevice url.Service
	if *mySQLConnectionString != "-" {
		sqlDB, err := db.InitMySQL(*mySQLConnectionString)
		if err != nil {
			panic(err)
		}
		urlSevice, err = db.NewURLInterfaceMySQL(sqlDB)
		if err != nil {
			panic(err)
		}
	} else if *redisAddressstring != "-" {
		var err error
		redisDB := db.InitRedis(*redisAddressstring)
		urlSevice, err = db.NewURLInterfaceRedis(redisDB)
		if err != nil {
			panic(err)
		}
	} else {
		panic("mysqlURL, redisURL both can't be null")
	}

	shortner := server.Shortner(*templatePath, urlSevice)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/", shortner.R)

	http.ListenAndServe(":"+strconv.Itoa(*port), r)
}
