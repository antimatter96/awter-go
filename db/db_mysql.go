// Package db Contains all methods used by the other functions
package db

import (
	"database/sql"

	url "github.com/antimatter96/awter-go/db/url"

	// This exposes mysql connector
	_ "github.com/go-sql-driver/mysql"
)

// Init is called by main since it requires
func InitMySQL(DBConnectionString string) (*sql.DB, error) {
	sqlDB, err := sql.Open("mysql", DBConnectionString)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(3)
	if err != nil {
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		sqlDB.Close()
		return nil, err
	}
	return sqlDB, nil
}

func NewURLInterfaceMySQL(sqlDB *sql.DB) (url.Service, error) {
	urlService := url.UrlsDb{DB: sqlDB}
	err := urlService.Init()
	return &urlService, err
}
