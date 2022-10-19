package main

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var a *application
var src = rand.NewSource(time.Now().UnixNano())

func NewHelpers(app *application) {
	a = app
}

// IsAuthenticated returns true if a user is authenticated
func IsAuthenticated(r *http.Request) bool {
	existsUserID := a.Session.Exists(r.Context(), "userID")
	existsToken := a.Session.Exists(r.Context(), "token")
	return existsUserID && existsToken
}

// RandomString returns a random string of letters of length n
func RandomString(n int) string {
	b := make([]byte, n)

	for i, theCache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			theCache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(theCache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		theCache >>= letterIdxBits
		remain--
	}

	return string(b)
}
