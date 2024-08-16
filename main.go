package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type templates struct {
	*template.Template
}

func my_init() {
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

func createTable() bool {
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

func createTable_account() bool {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS account (
        u_id SERIAL PRIMARY KEY,
        u_name VARCHAR(255),
		u_pass VARCHAR(255)
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

func populateTable_account(user string, pass string) {
	insertSQL := `INSERT INTO account (u_name, u_pass) VALUES ($1, $2)`

	// typically hash pw here
	hash_pass, err := hashPassword(pass)
	if err != nil {
		log.Printf("Unable to hash pass: %v", err)
	}
	_, err = db.Exec(insertSQL, user, hash_pass)
	if err != nil {
		log.Printf("Unable to create user: %v", err)
	}
	fmt.Println("Account created successfully!")
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

	// add log in table and root user
	if createTable_account() {
		populateTable_account("admin", "admin")
	}
}

// retreive for use in post methods
func getRandomRData() string {
	// Get the total number of entries
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM response").Scan(&count)
	if err != nil {
		return ""
	}

	// Generate a random ID
	randomID := rand.Intn(count) + 1

	// Fetch the r_data associated with the random ID
	var rData string
	err = db.QueryRow("SELECT r_data FROM response WHERE r_id = $1", randomID).Scan(&rData)
	if err != nil {
		return ""
	}

	return rData
}

func verifyLogin(u string, p string) bool {
	// Variables to store the retrieved data
	var dbPassword string

	// Query the database for the user's password
	query := `SELECT u_pass FROM account WHERE u_name=$1`
	err := db.QueryRow(query, u).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no user is found
			fmt.Printf("error")
			return false
		}
	}
	// Compare the stored plain text password with the provided password
	if !comparePasswords(dbPassword,p){
		return false
	}
	

	// If passwords match, handle successful login (e.g., create session)
	return true
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func comparePasswords(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// For rendering in later functions
func (t templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ExecuteTemplate(w, name, data)
}

func initTemplates() templates {
	t := template.New("")
	// Functions mapped to each template
	t.Funcs(template.FuncMap{})

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
	my_init()
	connDb()
	defer db.Close()

	// WEB INIT
	e := echo.New()
	e.Renderer = initTemplates()
	e.Use(middleware.Gzip(), middleware.Secure())

	e.GET("/", root)
	e.GET("/dash", dash)
	e.POST("/query", query)
	e.GET("/answer", answer)
	e.POST("/login", login)
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
		// change user spot to current user
		"user": "USERNAME",
		"title": "ChatGSC",
		"link": "/",
	})
}

// send users own words back
func query(c echo.Context) error {
	// validation and error check
	c.Request().ParseForm()
	unv_input := c.FormValue("user_txt")

	return c.Render(200, "chat.html", map[string]interface{}{
		// change to current user
		"user": "USERNAME",
		"q":    unv_input,
	})
}

func answer(c echo.Context) error {

	// gen rand data
	respo := getRandomRData()
	return c.Render(200, "answer.html", map[string]interface{}{
		"a": respo,
	})
}

// login
func login(c echo.Context) error {
	// validation and error check
	c.Request().ParseForm()
	username_input := c.FormValue("username")
	password_input := c.FormValue("password")


	// if strings are valid
	// call verify func
	if !verifyLogin(username_input, password_input) {
		return c.Render(200, "index.html", map[string]interface{}{
			"err": "INCORRECT",
		})
	}
	return c.Redirect(302, "/dash")
}
