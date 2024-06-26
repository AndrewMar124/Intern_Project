package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"log"
	"html/template"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}
	fmt.Println("Successfully connected to the database")
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", helloWorldHandler)
	
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
