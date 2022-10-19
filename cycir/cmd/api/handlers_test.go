package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name   string
	url    string
	method string
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
	{"one user", "/admin/user/{id}", "GET"},
	{"all hosts", "/admin/host/all", "GET"},
	{"one host", "/admin/host/{id}", "GET"},
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

var loginPostTests = []struct {
	name                string
	email               string
	remember            string
	expectedStatusCode  int
	expectedCookieValue string
	expectedHTML        string
	expectedLocation    string
}{
	{
		"with-remember-me",
		"admin@example.com",
		"remember",
		http.StatusSeeOther,
		"__gowatcher_remember=0",
		"",
		"/admin/overview",
	},
	{
		"without-remember-me",
		"jack@nimble.com",
		"",
		http.StatusSeeOther,
		"",
		"",
		"/admin/overview",
	},
}

func TestPostLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginPostTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("remember", e.remember)

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.CheckAuthentication)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedCookieValue != "" {
			// get the URL from test
			actualHeader := rr.Result().Header
			cookie := actualHeader.Get("Set-Cookie")
			cookieValue := strings.Split(cookie, "|")[0]
			if cookieValue != e.expectedCookieValue {
				t.Errorf("failed %s: expected cookie value %s, but got value %s", e.name, e.expectedCookieValue, cookieValue)
			}
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
