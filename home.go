package main

import (
    "html/template"
    "net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    tmpl.Execute(w, map[string]interface{}{
        "IsLoggedIn": session.Values["authenticated"] == true,
        "Username":   session.Values["username"],
    })
}