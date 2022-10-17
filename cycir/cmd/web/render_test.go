package main

import (
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}	

	testSession.Put(r.Context(), "flash", "123")

	result := testApp.DefaultData(td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var ww myWriter
	
	// add template functions
	err = testApp.RenderPage(&ww, r, "dashboard", nil, nil)

	if err != nil {
		t.Error("error writing template to browser")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = testSession.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

