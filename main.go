package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
	"os"
	_ "github.com/lib/pq"
)


var db *sql.DB
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func initDB() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Error opening database: %q", err)
	}
	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %q", err)
	}
	fmt.Println("Successfully connected to the database")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl.Execute(w, map[string]interface{}{
    })
}


func main() {
	// initialize database
	initDB()
	defer db.Close()

	http.HandleFunc("/", homeHandler)
	
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
