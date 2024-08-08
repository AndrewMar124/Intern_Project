package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB

type templates struct {
	*template.Template
}

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("Error opening database: %q", err)
	}
	fmt.Println("Successfully opened db!!")
	attemptDBAccess()
}

func attemptDBAccess() bool {
	err := db.Ping()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return false
	} else {
		log.Println("Successfully connected to the database!")
		return true
	}
}

func createTable() bool{
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS response (
        r_id SERIAL PRIMARY KEY,
        r_data VARCHAR(255)
    );`
    _, err := db.Exec(createTableSQL)
    if err != nil {
        log.Printf("Unable to create table: %v", err)
		return false
    }
    fmt.Println("Table created successfully!")
	return true
}

func populateTable(numRows int) {
    insertSQL := `INSERT INTO response (r_data) VALUES ($1)`

    for i := 0; i < numRows; i++ {
        sentence := generateRandomSentence()
        _, err := db.Exec(insertSQL, sentence)
        if err != nil {
            log.Fatalf("Unable to insert row: %v", err)
        }
    }
    fmt.Println("Table populated successfully!")
}

func generateRandomSentence() string {
    words := []string{"Lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit", "sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore", "magna", "aliqua"}
    sentenceLength := rand.Intn(10) + 5 // Generate a random sentence length between 5 and 15 words

    sentence := ""
    for i := 0; i < sentenceLength; i++ {
        sentence += words[rand.Intn(len(words))] + " "
    }
    return sentence
}

func connDb() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	// access db
	for range ticker.C {
		if attemptDBAccess() {
			break
		}
	}
	// attempt to populate the db
	if createTable() {
		populateTable(15)
	}
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
	// DB connect
	connDb()
	defer db.Close()

	

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
		"a":    slice[rand.Intn(6)],
	})
}

// toString converts any value to string
// functions that return a string are automatically escaped by html/template
func toString(v interface{}) string {
	return fmt.Sprint(v)
}
