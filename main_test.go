package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHandle(t *testing.T) {
	// t.Errorf("my test failed :(")
}

func TestDash(t *testing.T) {
	// Create an Echo instance
	e := echo.New()
	e.Renderer = initTemplates()

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/dash", nil)

	// Create a new HTTP recorder to capture the response
	rec := httptest.NewRecorder()

	// Create an Echo context
	c := e.NewContext(req, rec)
	err := dash(c)
	if err != nil {
		t.Error(err)
	}
    
	if rec.Code != http.StatusOK {
		t.Error("page not found")
	}
	if !strings.Contains(rec.Body.String(), "ChatGSC") {
		t.Error("missing data")
	}
	// Call the handler
	/*
		if assert.NoError(t, dash(c)) {
			// Check the status code
			assert.Equal(t, http.StatusOK, rec.Code)

			// Check the rendered template name
			assert.Contains(t, rec.Body.String(), "dash.html")

			// Check that the title is rendered correctly
			assert.Contains(t, rec.Body.String(), "ChatGSC")

			// Check that the link is rendered correctly
			assert.Contains(t, rec.Body.String(), "/")
		}
	*/
}
