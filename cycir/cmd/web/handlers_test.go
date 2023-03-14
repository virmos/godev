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
	{"one user", "/admin/user/0", "GET"},
	{"one user", "/admin/user/1", "GET"},
	{"all hosts", "/admin/host/all", "GET"},
	{"one host", "/admin/host/0", "GET"},
	{"one host", "/admin/host/1", "GET"},
}

// TestHandlers tests all routes that don't require extra tests (gets)
func TestHandlers(t *testing.T) {
	routes := getRoutes()
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
var loginGetTests = []struct {
	name               string
	userID             string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"with-user-id",
		"1",
		http.StatusSeeOther,
		"",
		"/admin/overview",
	},
}

func TestGetLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginGetTests {
		// create request
		req, _ := http.NewRequest("GET", "/", nil)
		ctx := getCtx(req)

		if e.userID != "" {
			testSession.Put(ctx, "userID", e.userID)
		}

		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.LoginScreen)
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

// loginPostTests is the data for the Login Post handler tests
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

// logoutTests is the data for the Logout handler tests
var logoutTests = []struct {
	name               string
	remember           string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"with-remember-me",
		"remember",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"without-remember-me",
		"",
		http.StatusSeeOther,
		"",
		"/",
	},
}

func TestLogout(t *testing.T) {
	// range through all tests
	for _, e := range logoutTests {
		// create request
		req, _ := http.NewRequest("GET", "/user/logout", nil)
		ctx := getCtx(req)

		if e.remember != "" {
			req.AddCookie(&http.Cookie{
				Name:     "__gowatcher_remember",
				Value:    "1|2",
				Domain:   "localhost",
				Path:     "/",
				MaxAge:   0,
				HttpOnly: true,
			})
		}

		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.Logout)
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

// saveInCacheTests is the data for the Login handler tests
var saveInCacheTests = []struct {
	name               string
	key                string
	value              string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"save in cache",
		"key",
		"value",
		http.StatusOK,
		"",
	},
}

func TestSaveInCache(t *testing.T) {
	// range through all tests
	for _, e := range saveInCacheTests {
		postedData := url.Values{}
		postedData.Add("key", e.key)
		postedData.Add("value", e.value)

		// create request
		req, _ := http.NewRequest("POST", "/save-in-cache", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.SaveInCache)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
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

// getFromCacheTests is the data for the Login handler tests
var getFromCacheTests = []struct {
	name               string
	key                string
	value              string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"save in cache",
		"key",
		"value",
		http.StatusOK,
		"",
	},
}

func TestGetFromCache(t *testing.T) {
	// range through all tests
	for _, e := range getFromCacheTests {
		postedData := url.Values{}
		postedData.Add("key", e.key)
		postedData.Add("value", e.value)

		// create request
		req, _ := http.NewRequest("POST", "/get-from-cache", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.GetFromCache)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
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

// deleteFromCacheTests is the data for the Login handler tests
var deleteFromCacheTests = []struct {
	name               string
	key                string
	value              string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"save in cache",
		"key",
		"value",
		http.StatusOK,
		"",
	},
}

func TestDeleteFromCache(t *testing.T) {
	// range through all tests
	for _, e := range deleteFromCacheTests {
		postedData := url.Values{}
		postedData.Add("key", e.key)
		postedData.Add("value", e.value)

		// create request
		req, _ := http.NewRequest("POST", "/delete-from-cache", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.DeleteFromCache)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
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

// emptyCacheTests is the data for the Login handler tests
var emptyCacheTests = []struct {
	name               string
	key                string
	value              string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"save in cache",
		"key",
		"value",
		http.StatusOK,
		"",
	},
}

func TestEmptyCache(t *testing.T) {
	// range through all tests
	for _, e := range emptyCacheTests {
		postedData := url.Values{}
		postedData.Add("key", e.key)
		postedData.Add("value", e.value)

		// create request
		req, _ := http.NewRequest("POST", "/empty-cache", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.EmptyCache)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
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
