package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestDash(t *testing.T) {
	t.Run("render template", func(t *testing.T) {
		ptitle := "ChatGSC"
		rec := createContext(t, "/dash")
		assertCorrectResponse(t, rec, ptitle)
	})
	
}

func assertCorrectResponse(t testing.TB, rec *httptest.ResponseRecorder, title string) {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Errorf("Status expected %q, status received %q", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), title) {
		t.Errorf("Missing page content: %q", title)
	}
}

func createContext(t testing.TB, page string) *httptest.ResponseRecorder{
	// Create an Echo instance
	e := echo.New()
	e.Renderer = initTemplates()

	// req + rec
	req := httptest.NewRequest(http.MethodGet, page, nil)
	// Create a new HTTP recorder to capture the response
	rec := httptest.NewRecorder()

	// make context
	c := e.NewContext(req, rec)
	err := dash(c)
	if err != nil {
		t.Error(err)
	}

	return rec
}
