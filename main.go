package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"strings"
	"database/sql"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB

type templates struct {
	*template.Template
}

type Response struct {
	ID int
	Rdata string
}

func init(){
	// DB INIT
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func getResponse() ([]Response) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil
	}
	defer rows.Close()
	
	users := []Response{}
	for rows.Next() {
		var user Response
		err := rows.Scan(&user.ID, &user.Rdata)
		if err != nil {
			return nil
		}
		users = append(users, user)
	}
	return users
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
		"title": "Root",
		"test":  "Hello, world!",
		"slice": []int{1, 2, 3},
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
	//response := getResponse()

	return c.Render(200, "chat.html", map[string]interface{}{
		"user": "USERNAME",
		"q":    unv_input,
		//"a": response[0],
	})
}

// toString converts any value to string
// functions that return a string are automatically escaped by html/template
func toString(v interface{}) string {
	return fmt.Sprint(v)
}
