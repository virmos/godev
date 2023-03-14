package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

var authTests = []struct {
	name                string
	authorization				string
	expectedStatusCode  int
	expectedHTML        string
	expectedLocation    string
}{
	{
		"middleware auth",
		"Bearer 12345678912345678912345678",
		http.StatusOK,
		"",
		"",
	},
	{
		"middleware auth",
		"",
		http.StatusUnauthorized,
		"",
		"",
	},
}

func TestAuth(t *testing.T) {
	for _, e := range authTests {
		req, _ := http.NewRequest("GET", "/", nil)

		// set the header
		req.Header.Set("Content-Type", "text/html")
		req.Header.Set("Authorization", e.authorization)

		var myH myHandler

		h := testApp.Auth(&myH)

		switch v := h.(type) {
		case http.Handler:
			rr := httptest.NewRecorder()
			v.ServeHTTP(rr, req)

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

		default:
			t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
		}
	}
}
