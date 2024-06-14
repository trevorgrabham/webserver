package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func ConnectDb() (*sql.DB, error) {
	cfg := mysql.Config{
		User:		os.Getenv("DBUSER"),
		Passwd:	os.Getenv("DBPASS"),
		Net:		"tcp",
		Addr:		"127.0.0.1:3306",
		DBName:	"testing_db",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("Opening mysql: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Connecting to db: %v", err)
	}
	return db, nil
}
