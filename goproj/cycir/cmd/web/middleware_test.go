package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestNoPanic(t *testing.T) {
	var myH myHandler

	h := RecoverPanic(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

var checkRememberMeTests = []struct {
	name           string
	remember       string
	userID         interface{}
	hash           string
	token          interface{}
	expectedUserID interface{}
	expectedToken  interface{}
}{
	{
		"credentials + no remember",
		"",
		"1",
		"xyz",
		"abc",
		"1",
		"abc",
	},
	{
		"credentials + remember",
		"remember",
		"1",
		"xyz", // valid hash
		"abc",
		"1",
		"abc",
	},
	{
		"credentials + remember",
		"remember",
		"1",
		"", // invalid hash
		"abc",
		nil,
		nil,
	},
	{
		"invalid credentials + no remember",
		"",
		"1",
		"xyz",
		"abc",
		"1",
		"abc",
	},
	{
		"invalid credentials + remember",
		"remember",
		"",
		"xyz", // valid hash
		"abc",
		"0",
		"abc",
	},
	{
		"invalid credentials + remember",
		"remember",
		"0",
		"xyz", // invalid hash
		"abc",
		nil,
		nil,
	},
}

func TestCheckRemember(t *testing.T) {
	for _, e := range checkRememberMeTests {
		req, _ := http.NewRequest("GET", "/", nil)
		ctx := getCtx(req)

		if e.remember == "remember" {
			req.AddCookie(&http.Cookie{
				Name:     "__gowatcher_remember",
				Value:    fmt.Sprintf("%s|%s", e.userID, e.hash),
				Domain:   "localhost",
				Path:     "/",
				MaxAge:   0,
				HttpOnly: true,
			})
		}
		if e.name == "credentials + no remember" || e.name == "credentials + remember" { // so isAuthenticated returns true
			testSession.Put(ctx, "usedID", e.userID)
		}
		testSession.Put(ctx, "token", e.token)

		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "text/html")

		var myH myHandler

		h := testApp.CheckRemember(&myH)

		switch v := h.(type) {
		case http.Handler:
			rr := httptest.NewRecorder()
			v.ServeHTTP(rr, req)
		default:
			t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
		}
	}
}
