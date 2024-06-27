package main

import (
	"net/http"
	"github.com/gorilla/sessions"
)

// Global variable for the session store
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// login page *fix with secure login / DB integration
func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")

        // Here, add your logic to authenticate the user
        if username == "user" && password == "pass" {
            session, _ := store.Get(r, "session-name")
            session.Values["authenticated"] = true
            session.Values["username"] = username
            session.Save(r, w)
            http.Redirect(w, r, "/", http.StatusFound)
            return
        }
    }
    http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    session.Values["authenticated"] = false
    session.Values["username"] = ""
    session.Save(r, w)
    http.Redirect(w, r, "/", http.StatusFound)
}
