package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/antimatter96/awter-go/customcrypto"
	"github.com/antimatter96/awter-go/db"
	"github.com/antimatter96/awter-go/db/url"
	"github.com/antimatter96/awter-go/server"
)

var build string

func main() {
	fmt.Printf("Starting build: %s\n", build)
	var port = flag.Int("port", 8080, "port")
	var templatePath = flag.String("template", "./template", "the template directory")
	var mySQLConnectionString = flag.String("mysqlURL", "user:password@/name?parseTime=true", "MySQL connection string")
	var redisAddressstring = flag.String("redisURL", "0.1:6379", "redis connection string")

	flag.Parse()

	fmt.Printf("Build params:\nPort: %d\nTemplate: %s\nMySQL: %s\nRedis: %s\n%s\n", *port, *templatePath, *mySQLConnectionString, *redisAddressstring, "HELLO")

	var urlSevice url.Service
	if *mySQLConnectionString != "" {
		sqlDB, err := db.InitMySQL(*mySQLConnectionString)
		if err != nil {
			panic(err)
		}
		urlSevice, err = db.NewURLInterfaceMySQL(sqlDB)
		if err != nil {
			panic(err)
		}
	} else if *redisAddressstring != "" {
		var err error
		redisDB := db.InitRedis(*redisAddressstring)
		urlSevice, err = db.NewURLInterfaceRedis(redisDB)
		if err != nil {
			panic(err)
		}
	} else {
		panic("mysqlURL, redisURL both can't be null")
	}

	naclScrypt := &customcrypto.NaclScrypt{}
	passwordChecker := &customcrypto.Bcrypt{Cost: 12}

	shortner := server.Shortner(*templatePath, urlSevice, naclScrypt, passwordChecker)

	logger := newLogger()

	r := newRouter(shortner, logger)

	http.ListenAndServe(":"+strconv.Itoa(*port), r)
}

func newRouter(shortner *server.Server, logger zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(hlog.NewHandler(logger))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/", shortner.R)

	return r
}

func newLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return log
}
