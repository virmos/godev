package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"strings"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
}{
	{"login page", "/", "GET"},
	{"dashboard page", "/admin/overview", "GET"},
	{"events page", "/admin/events", "GET"},
	{"schedule page", "/admin/schedule", "GET"},
	{"settings page", "/admin/settings", "GET"},
	{"all healthy", "/admin/all-healthy", "GET"},
	{"all warnings", "/admin/all-warning", "GET"},
	{"all problems", "/admin/all-problems", "GET"},
	{"all pendings", "/admin/all-pending", "GET"},
	{"all users", "/admin/users", "GET"},
	{"one user", "/admin/users/{id}", "GET"},	
	{"all hosts", "/admin/hosts", "GET"},
	{"one host", "/admin/hosts/{id}", "GET"},
}

// TestHandlers tests all routes that don't require extra tests (gets)
func TestHandlers(t *testing.T) {
	routes := testApp.routes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		_, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Errorf("for %s, error rendering", e.name)
		}
	}
}


// loginTests is the data for the Login handler tests
var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"admin@example.com",
		http.StatusSeeOther,
		"",
		"/admin/overview",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}


// gets the context
func getCtx(req *http.Request) context.Context {
	ctx, err := testApp.Session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
