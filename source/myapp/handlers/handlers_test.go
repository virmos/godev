package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHome(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/")
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("for home page, expected status 200 but got %d", resp.StatusCode)
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(bodyText), "awesome") {
		cel.TakeScreenShot(ts.URL + "/", "HomeTest", 1500, 1000)
		t.Error("did not find submarine")
	}
}

func TestClicker(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	page := cel.FetchPage(ts.URL + "/tester")
	outputElement := cel.SelectElementByID(page, "output")
	button := cel.SelectElementByID(page, "clicker")

	testHTML, _ := outputElement.HTML()
	if strings.Contains(testHTML, "Clicked the button") {
		t.Errorf("found text that should not be there")
	}

	button.MustClick()
	testHTML, _ = outputElement.HTML()
	if !strings.Contains(testHTML, "Clicked the button") {
		t.Errorf("did not find text that should be there")
	}
}

func TestHome2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	cel.Session.Put(ctx, "test_key", "Hello, world.")

	h := http.HandlerFunc(testHandlers.Home)
	h.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("returned wrong response code; got %d but expected 200", rr.Code)
	}

	if cel.Session.GetString(ctx, "test_key") != "Hello, world." {
		t.Error("did not get correct value from session")
	}
}