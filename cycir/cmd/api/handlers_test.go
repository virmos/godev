package main

import (
	"bytes"
	"context"
	"cycir/internal/elastics"
	"cycir/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type respData struct {
	Error          bool   `json:"error"`
	Message        string `json:"message"`
	Redirect       bool   `json:"redirect"`
	RedirectStatus int    `json:"status"`
	Route          string `json:"route"`
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
		req, _ := http.NewRequest("POST", "api/authenticate", bytes.NewReader(body))

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
	expectedLocation   string
}{
	{
		"valid user",
		"1",
		http.StatusOK,
		"/admin/users",
	},
}

func TestDeleteUser(t *testing.T) {
	// range through all tests
	for _, e := range deleteUserTests {
		postBody := map[string]interface{}{}

		body, _ := json.Marshal(postBody)

		// create request
		req, _ := http.NewRequest("POST", fmt.Sprintf("api/admin/user/delete/%s", e.id), bytes.NewReader(body))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", e.id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.DeleteUser)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		// checking for expected values in HTML
		if e.expectedLocation != "" {
			// read the response body into a string
			var out respData
			dec := json.NewDecoder(rr.Body)
			dec.Decode(&out)
			if out.Route != e.expectedLocation {
				t.Errorf("failed %s: expected code %s, but got %s", e.name, e.expectedLocation, out.Route)
			}
		}
	}
}

var postSettingTests = []struct {
	name               string
	action             string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"redirect",
		"1",
		http.StatusOK,
		"",
		"/admin/overview",
	},
	{
		"not redirect",
		"0",
		http.StatusOK,
		"",
		"/admin/settings",
	},
}

func TestPostSettings(t *testing.T) {
	// range through all tests
	for _, e := range postSettingTests {
		var setting models.Setting
		setting.Action = e.action

		body, _ := json.Marshal(setting)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/settings", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.PostSettings)
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

		if e.expectedLocation != "" {
			// read the response body into a string
			var out respData
			dec := json.NewDecoder(rr.Body)
			dec.Decode(&out)
			if out.Route != e.expectedLocation {
				t.Errorf("failed %s: expected code %s, but got %s", e.name, e.expectedLocation, out.Route)
			}
		}
	}
}

var postHostTests = []struct {
	name               string
	id                 string
	action             string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid",
		"1",
		"1",
		http.StatusOK,
		"",
		"/admin/host/all",
	},
	{
		"invalid",
		"0",
		"0",
		http.StatusOK,
		"",
		"/admin/host/1",
	},
}

func TestPostHost(t *testing.T) {
	// range through all tests
	for _, e := range postHostTests {
		var h models.Host
		h.Action = e.action

		body, _ := json.Marshal(h)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/host/1", bytes.NewReader(body))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", e.id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		// set the header
		req.Header.Set("Content-Type", "application/json")
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
		if e.expectedLocation != "" {
			// read the response body into a string
			var out respData
			dec := json.NewDecoder(rr.Body)
			dec.Decode(&out)
			if out.Route != e.expectedLocation {
				t.Errorf("failed %s: expected code %s, but got %s", e.name, e.expectedLocation, out.Route)
			}
		}
	}
}

var deleteHostTests = []struct {
	name               string
	id                 string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid",
		"1",
		http.StatusOK,
		"",
		"",
	},
	{
		"invalid",
		"0",
		http.StatusOK,
		"",
		"",
	},
}

func TestDeleteHost(t *testing.T) {
	// range through all tests
	for _, e := range deleteHostTests {
		var u models.User

		body, _ := json.Marshal(u)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/host/delete/1", bytes.NewReader(body))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", e.id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.DeleteHost)
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
		if e.expectedLocation != "" {
			// read the response body into a string
			var out respData
			dec := json.NewDecoder(rr.Body)
			dec.Decode(&out)
			if out.Route != e.expectedLocation {
				t.Errorf("failed %s: expected code %s, but got %s", e.name, e.expectedLocation, out.Route)
			}
		}
	}
}

var postOneUserTests = []struct {
	name               string
	id                 string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid",
		"1",
		http.StatusOK,
		"",
		"/admin/users",
	},
	{
		"invalid",
		"0",
		http.StatusOK,
		"",
		"/admin/users",
	},
}

func TestPostOneUser(t *testing.T) {
	// range through all tests
	for _, e := range postOneUserTests {
		var u models.User

		body, _ := json.Marshal(u)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/user/1", bytes.NewReader(body))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", e.id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		// set the header
		req.Header.Set("Content-Type", "application/json")
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
		if e.expectedLocation != "" {
			// read the response body into a string
			var out respData
			dec := json.NewDecoder(rr.Body)
			dec.Decode(&out)
			if out.Route != e.expectedLocation {
				t.Errorf("failed %s: expected code %s, but got %s", e.name, e.expectedLocation, out.Route)
			}
		}
	}
}

var setSystemPrefTests = []struct {
	name               string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		http.StatusOK,
		"",
	},
}

func TestSetSystemPref(t *testing.T) {
	// range through all tests
	for _, e := range setSystemPrefTests {
		var payload struct {
			PrefName  string `json:"pref_name"`
			PrefValue string `json:"pref_value"`
		}
		testApp.PreferenceMap = make(map[string]string)

		body, _ := json.Marshal(payload)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/preference/ajax/set-system-pref", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
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
	active             int
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		1,
		http.StatusOK,
		"",
	},
	{
		"invalid",
		0,
		http.StatusOK,
		"",
	},
}

func TestToggleServiceForHost(t *testing.T) {
	// range through all tests
	for _, e := range toggleServiceForHostTests {
		var payload struct {
			HostID    int `json:"host_id"`
			ServiceID int `json:"service_id"`
			Active    int `json:"active"`
		}
		payload.Active = e.active

		body, _ := json.Marshal(payload)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/host/ajax/toggle-service", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
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
	enabled            string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"1",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"0",
		http.StatusOK,
		"",
	},
}

func TestToggleMonitoring(t *testing.T) {
	// range through all tests
	for _, e := range toggleMonitoringTests {
		var payload struct {
			Enabled string `json:"enabled"`
		}
		payload.Enabled = e.enabled

		body, _ := json.Marshal(payload)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/preference/ajax/toggle-monitoring", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
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

var sendRangeUptimeReportTests = []struct {
	name               string
	host_name          string
	start_date         string
	end_date           string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"Google",
		"2022-10-01",
		"2022-10-29",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"Google",
		"2022-10-20",
		"2022-10-20",
		http.StatusOK,
		"",
	},
}

func TestSendRangeUptimeReport(t *testing.T) {
	// range through all tests
	for _, e := range sendRangeUptimeReportTests {
		var payload struct {
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
			HostName  string `json:"host_name"`
			CSRF      string `json:"csrf_token"`
		}
		payload.HostName = e.host_name
		payload.StartDate = e.start_date
		payload.EndDate = e.end_date

		body, _ := json.Marshal(payload)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/host/send-range-uptime-report", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.SendRangeUptimeReport)
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

var sendRangeUptimeReportCachedTests = []struct {
	name               string
	host_name          string
	start_date         string
	end_date           string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		"valid",
		"Google",
		"2022-10-01",
		"2022-10-29",
		http.StatusOK,
		"",
	},
	{
		"invalid",
		"Google",
		"2022-10-20",
		"2022-10-20",
		http.StatusOK,
		"",
	},
}

func TestSendRangeUptimeReportCached(t *testing.T) {
	// range through all tests
	for _, e := range sendRangeUptimeReportCachedTests {
		var payload struct {
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
			HostName  string `json:"host_name"`
			Histogram string `json:"histogram"`
			Count     string `json:"count"`
			CSRF      string `json:"csrf_token"`
		}
		payload.HostName = e.host_name
		payload.StartDate = e.start_date
		payload.EndDate = e.end_date
		payload.Histogram = ""
		payload.Count = ""

		body, _ := json.Marshal(payload)

		// create request
		req, _ := http.NewRequest("POST", "api/admin/host/send-range-uptime-report-cached", bytes.NewReader(body))

		// set the header
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(testApp.SendRangeUptimeReportCached)
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

var helperRangeUptimeReportTests = []struct {
	name                   string
	host_name              string
	hours_uptime_histogram []int
	hours_histogram        []int
	expectedHistogram      []string
	expectedCount          []string
}{
	{
		"valid",
		"Google",
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24},
		[]string{"0%", "50%", "66%", "75%", "80%", "83%", "85%", "87%", "88%", "90%", "90%", "91%", "92%", "92%", "93%", "93%", "94%", "94%", "94%", "95%", "95%", "95%", "95%", "95%"},
		[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24"},
	},
}

func TestHelperParseUptimeReports(t *testing.T) {
	for _, e := range helperRangeUptimeReportTests {
		uptimeReports := make(map[string]elastics.Report)
		reports := make(map[string]elastics.Report)

		newUptimeReport := elastics.Report{
			HoursHistogram: e.hours_uptime_histogram,
		}

		newReport := elastics.Report{
			HoursHistogram: e.hours_histogram,
		}

		uptimeReports["Google"] = newUptimeReport
		reports["Google"] = newReport

		results, _ := parseUptimeReports(uptimeReports, reports)
		for i, v := range results["Google"].Histogram {
			if v != e.expectedHistogram[i] {
				t.Errorf("failed %s: expected code %s, but got %s", e.name, e.expectedHistogram[i], v)
			}
		}

		for i, v := range results["Google"].Count {
			if v != e.expectedCount[i] {
				t.Errorf("failed %s: expected code %s, but got %s", e.name, e.expectedHistogram[i], v)
			}
		}
	}
}

var helperRangeUptimeRangeReportTests = []struct {
	name                  string
	host_name             string
	days_uptime_histogram []int
	days_histogram        []int
	expectedHistogram     []string
	expectedCount         []string
}{
	{
		"valid",
		"Google",
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		[]string{"0%","50%","66%","75%","80%","83%","85%","87%","88%","90%","90%","91%","92%","92%","93%","93%","94%","94%","94%","95%","95%","95%","95%","95%","96%","96%","96%","96%","96%","96%","96%"},
		[]string{"1","2","3","4","5","6","7","8","9","10","11","12","13","14","15","16","17","18","19","20","21","22","23","24","25","26","27","28","29","30","31"},
	},
}

func TestHelperParseUptimeRangeReports(t *testing.T) {
	for _, e := range helperRangeUptimeRangeReportTests {
		uptimeReports := make(map[string]elastics.Report)
		reports := make(map[string]elastics.Report)

		newUptimeReport := elastics.Report{
			DaysHistogram: e.days_uptime_histogram,
		}

		newReport := elastics.Report{
			DaysHistogram: e.days_histogram,
		}

		uptimeReports["Google"] = newUptimeReport
		reports["Google"] = newReport

		results, _ := parseUptimeRangeReports(uptimeReports, reports)
		for i, v := range results["Google"].Histogram {
			if v != e.expectedHistogram[i] {
				t.Errorf("failed %s: expected code %s, but got %s", e.name, e.expectedHistogram[i], v)
			}
		}

		for i, v := range results["Google"].Count {
			if v != e.expectedCount[i] {
				t.Errorf("failed %s: expected code %s, but got %s", e.name, e.expectedHistogram[i], v)
			}
		}
	}
}
