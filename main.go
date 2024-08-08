package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"strings"
	"math/rand"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type templates struct {
	*template.Template
}

// For rendering in later functions
func (t templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ExecuteTemplate(w, name, data)
}

func initTemplates() templates {
	t := template.New("")
	// Functions mapped to each template
	t.Funcs(template.FuncMap{
		"toString": toString,
	})

	// parse all html files in the templates directory
	t, err := t.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// also parse files in subdirectories of templates
	//t, err = t.ParseGlob("templates/*/*.html")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// all templates created
	return templates{t}
}

func main() {
	// WEB INIT
	e := echo.New()
	e.Renderer = initTemplates()
	e.Use(middleware.Gzip(), middleware.Secure())

	e.GET("/", root)
	e.GET("/dash", dash)
	e.POST("/query", query)
	e.Static("/dist", "./dist")
	e.Start(":3000")
}

// index.html file
func root(c echo.Context) error {
	return c.Render(200, "index.html", map[string]interface{}{
		"title": "Home",
		"test":  "Welcome...",
		"link":  "/dash",
	})
}

func dash(c echo.Context) error {
	return c.Render(200, "dash.html", map[string]interface{}{
		"title": "ChatGSC",
		// @todo change this to username from db
		//"user": "USERNAME",
		"link": "/",
	})
}

// send users own words back
func query(c echo.Context) error {
	// validation and error check
	c.Request().ParseForm()
	unv_input := c.FormValue("user_txt")
	// validate input
	if strings.Contains(unv_input, "<") ||
		strings.Contains(unv_input, ">") {
		unv_input = "ERROR - INVALID INPUT"
		
	}

	// send response to AI... aka DB
	// responses := getResponse()
	slice := []string{"Hi!", "Sorry I can't help with that...", 
	"This information can be found in GSC policy number 233",
	 "Holo", "OK!", "NO!"}

	return c.Render(200, "chat.html", map[string]interface{}{
		"user": "USERNAME",
		"q":    unv_input,
		"a": slice[rand.Intn(6)],
	})
}

// toString converts any value to string
// functions that return a string are automatically escaped by html/template
func toString(v interface{}) string {
	return fmt.Sprint(v)
}
