package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	mux := app.routes()
	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing; test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, type is %T", v))
	}
}
