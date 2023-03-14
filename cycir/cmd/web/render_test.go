package main

import (
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
	err = testApp.RenderPage(&ww, r, "login", nil, nil)

	if err != nil {
		t.Error("error writing template to browser")
	}
}
