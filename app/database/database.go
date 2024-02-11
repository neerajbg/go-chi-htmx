package database

import (
	"database/sql"
	"log"
)

var DBConn *sql.DB

func ConnectDB() {
	dsn := "host=localhost port=5432 user=postgres password=neeraj dbname=chi-htmx-demo sslmode=disable"

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Println("Error in DB connection", err)
	}

	DBConn = db
	log.Println("Database connection successful.")

}
