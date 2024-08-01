package main

import (
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		"link":     link,
		"email": func() string {
			return "example@example.com"
		},
	})

	// parse all html files in the templates directory
	t, err := t.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// also parse files in subdirectories of templates
	t, err = t.ParseGlob("templates/*/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// all templates created
	return templates{t}
}

func main() {
	e := echo.New()

	e.Debug = true // TODO: this line should be removed in production
	// INIT templates func call
	e.Renderer = initTemplates()
	e.Use(middleware.Gzip(), middleware.Secure())
	// you should also add middleware.CSRF(), once you have forms

	e.GET("/", root)
	e.GET("/foo", foo)
	e.GET("/bar", bar)
	// fix
	e.GET("/dash", dash)
	e.Static("/dist", "./dist")
	e.Start(":3000")
}

// index.html file
func root(c echo.Context) error {
	return c.Render(200, "index.html", map[string]interface{}{
		"title": "Root",
		"test":  "Hello, world!",
		"slice": []int{1, 2, 3},
	})
}

func foo(c echo.Context) error {
	return c.Render(200, "foo.html", map[string]interface{}{
		"title": "Foo",
	})
}

func bar(c echo.Context) error {
	return c.Render(200, "bar.html", map[string]interface{}{
		"title": "Bar",
	})
}

func dash(c echo.Context) error {
	return c.Render(200, "dash.html", map[string]interface{}{
		"title": "ChatGSC",
		// @todo change this to username from db
		"user": "USERNAME",
	})
}

// toString converts any value to string
// functions that return a string are automatically escaped by html/template
func toString(v interface{}) string {
	return fmt.Sprint(v)
}

// link returns a styled "a" tag
// functions that return a template.HTML are not escaped, so all parameters need to be escaped to avoid xss
func link(location, name string) template.HTML {
	return escSprintf(`<a class="nav" href="%v">%v</a>`, location, name)
}

// escSprintf is like fmt.Sprintf but uses the escaped HTML equivalent of the args
func escSprintf(format string, args ...interface{}) template.HTML {
	for i, arg := range args {
		args[i] = template.HTMLEscapeString(fmt.Sprint(arg))
	}

	return template.HTML(fmt.Sprintf(format, args...))
}
