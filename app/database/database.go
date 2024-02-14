package database

import (
	"database/sql"
	"log"
	"os"
)

var DBConn *sql.DB

func ConnectDB() {

	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_USER := os.Getenv("DB_USER")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")

	dsn := "host=" + DB_HOST + " port=" + DB_PORT + " user=" + DB_USER + " password=" + DB_PASSWORD + " dbname=" + DB_NAME + " sslmode=disable"

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Println("Error in DB connection", err)
	}

	DBConn = db
	log.Println("Database connection successful.")

}
