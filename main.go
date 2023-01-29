package main

import (
	"./controllers"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

func main() {
	os.Remove("./tinder.db")

	// Open DB connection
	database, err := sql.Open("sqlite3", "./tinder.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Create DB Tables
	statement, _ := database.Prepare(
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, email TEXT, password TEXT, name TEXT, gender TEXT, age INTEGER, token TEXT, UNIQUE(email))")
	statement.Exec()

	statement, _ = database.Prepare(
		"CREATE TABLE IF NOT EXISTS swipes (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, potential_match_id INTEGER, preference INTEGER, match INTEGER, UNIQUE(user_id, potential_match_id))")
	statement.Exec()

	// MUX Router
	r := mux.NewRouter()

	// Create TinderService
	service := controllers.TinderService{
		Db: database,
	}

	// Routes
	r.HandleFunc("/login", service.Login).Methods("POST")
	r.HandleFunc("/profiles/{id:[0-9]+}", service.Profiles).Methods("POST")
	r.HandleFunc("/user/create", service.CreateUser).Methods("GET")
	r.HandleFunc("/swipe", service.Swipe).Methods("POST")

	// Run Web Server
	fmt.Printf("Starting the server on localhost:%d...\n", 3000)
	http.ListenAndServe(":3000", r)
}
