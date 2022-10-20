package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var createAuthedTokenTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"admin@example.com",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"jack@nimble.com",
		http.StatusUnauthorized,
		"",
	},
}

func TestCreateAuthedToken(t *testing.T) {
	// range through all tests
	for _, e := range createAuthedTokenTests {
		postBody := map[string]interface{}{
			"email":    e.email,
			"password": "",
		}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", "api/authenticated", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.CreateAuthToken)
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

var authenticationTests = []struct {
	name               string
	authorization      string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"Bearer 12345678912345678912345678",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"",
		http.StatusUnauthorized,
		"",
	},
}

func TestAuthentication(t *testing.T) {
	// range through all tests
	for _, e := range authenticationTests {
		postBody := map[string]interface{}{}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", "api/is-authenticated", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", e.authorization)
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.CheckAuthentication)
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

var deleteUserTests = []struct {
	name               string
	id                 string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid user",
		"1",
		http.StatusOK,
		"",
		"admin/users",
	},
	{
		"invalid user",
		"0",
		http.StatusBadRequest,
		"",
		"admin/users",
	},
}

func TestDeleteUser(t *testing.T) {
	// range through all tests
	for _, e := range deleteUserTests {
		postBody := map[string]interface{}{}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", fmt.Sprintf("api/admin/user/delete/%s", e.id), bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.DeleteUser)
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

var postSettingTests = []struct {
	name               string
	authorization      string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid",
		"Bearer 12345678912345678912345678",
		http.StatusOK,
		"",
		"admin/settings",
	},
	{
		"invalid",
		"",
		http.StatusUnauthorized,
		"",
		"",
	},
}

func TestPostSettings(t *testing.T) {
	// range through all tests
	for _, e := range postSettingTests {
		postBody := map[string]interface{}{}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/settings", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", e.authorization)
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.PostSettings)
		log.Println(testApp.repo)
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

var postHostTests = []struct {
	name               string
	authorization      string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"Bearer 12345678912345678912345678",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"",
		http.StatusUnauthorized,
		"",
	},
}

func TestPostHost(t *testing.T) {
	// range through all tests
	for _, e := range postHostTests {
		postBody := map[string]interface{}{}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", "api/host{id}", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", e.authorization)
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.PostHost)
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

var postOneUserTests = []struct {
	name               string
	authorization      string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"Bearer 12345678912345678912345678",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"",
		http.StatusUnauthorized,
		"",
	},
}

func TestPostOneUser(t *testing.T) {
	// range through all tests
	for _, e := range authenticationTests {
		postBody := map[string]interface{}{}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/user/{id}", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", e.authorization)
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.PostOneUser)
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

var setSystemPrefTests = []struct {
	name               string
	authorization      string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"Bearer 12345678912345678912345678",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"",
		http.StatusUnauthorized,
		"",
	},
}

func TestToggleSystemPref(t *testing.T) {
	// range through all tests
	for _, e := range setSystemPrefTests {
		postBody := map[string]interface{}{}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/preference/ajax/set-system-pref", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", e.authorization)
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.SetSystemPref)
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

var toggleServiceForHostTests = []struct {
	name               string
	authorization      string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"Bearer 12345678912345678912345678",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"",
		http.StatusUnauthorized,
		"",
	},
}

func TestToggleServiceForHost(t *testing.T) {
	// range through all tests
	for _, e := range toggleServiceForHostTests {
		postBody := map[string]interface{}{}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/host/ajax/toggle-service", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", e.authorization)
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.ToggleServiceForHost)
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

var toggleMonitoringTests = []struct {
	name               string
	authorization      string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"Bearer 12345678912345678912345678",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"",
		http.StatusUnauthorized,
		"",
	},
}

func TestToggleMonitoring(t *testing.T) {
	// range through all tests
	for _, e := range toggleMonitoringTests {
		postBody := map[string]interface{}{}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/preference/ajax/toggle-monitoring", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", e.authorization)
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.ToggleMonitoring)
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
