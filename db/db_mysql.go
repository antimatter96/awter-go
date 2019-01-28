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
var sqlDB *sql.DB

// Init is called by main since it requires
func InitMySQL() {
	DBConnectionString, _ := constants.Value("db-connection-string").(string)

	var err error
	sqlDB, err = sql.Open("mysql", DBConnectionString)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(3)
	if err != nil {
		panic(err.Error())
	}
	err = sqlDB.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func NewURLInterfaceMySQL() url.Service {
	urlService := url.UrlsDb{DB: sqlDB}
	err := urlService.Init()
	if err != nil {
		panic(err.Error())
	}
	return &urlService
}

func checkStatus() bool {
	err := sqlDB.Ping()
	if err != nil {
		return false
	}
	return true
}
