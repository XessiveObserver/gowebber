package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import postgres sql driver
)

// Loads enviroment variables from .env file
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// variable that holds database connection
var DB *sql.DB

var store *sessions.CookieStore

// Initialize database connection
func InitDB() {

	loadEnv()

	store = sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))
	// Load all enviroment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	DB = db
	fmt.Println("Connected to the database")
}
