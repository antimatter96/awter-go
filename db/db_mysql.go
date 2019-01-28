// Package db Contains all methods used by the other functions
package db

import (
	"database/sql"

	"../constants"
	url "./url"

	// This exposes mysql connector
	_ "github.com/go-sql-driver/mysql"
)

// The main db object
var db *sql.DB

// Init is called by main since it requires
func InitMySQL() {
	DBConnectionString, _ := constants.Value("db-connection-string").(string)

	var err error
	db, err = sql.Open("mysql", DBConnectionString)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(3)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func NewURLInterfaceMySQL() url.Service {
	urlService := url.UrlsDb{DB: db}
	err := urlService.Init()
	if err != nil {
		panic(err.Error())
	}
	return &urlService
}

func checkStatus() bool {
	err := db.Ping()
	if err != nil {
		return false
	}
	return true
}
