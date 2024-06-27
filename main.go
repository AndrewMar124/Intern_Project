package main

import (
	"fmt"
	"net/http"
)

func main() {
	// initialize database
	initDB()
	defer db.Close()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/logout", logoutHandler)
	
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
