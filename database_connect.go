//Include db functions here
package main

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/lib/pq"
)

var db *sql.DB

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
