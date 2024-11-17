package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func InitDataBase() (database *sql.DB) {
	// Database connection
	var db *sql.DB
	var err error
	cfg := mysql.Config{
		User:                 "user",
		Passwd:               "",
		Addr:                 "string",
		Net:                  "tcp",
		DBName:               "database",
		AllowNativePasswords: true,
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to the database")
	return db
}
