package main

import (
	"cycir/internal/models"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
)

var testApp *application

func TestMain(m *testing.M) {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Halifax")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	testApp = &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		repo:     models.NewTestRepository(),
	}

	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
