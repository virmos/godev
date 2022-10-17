package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		_, err := app.authenticateToken(r)
		if err != nil {
			app.invalidCredentials(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// NoSurf implements CSRF protection
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.ExemptPath("/pusher/auth")
	csrfHandler.ExemptPath("/pusher/hook")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   cfg.InProduction,
		SameSite: http.SameSiteStrictMode,
		Domain:   cfg.Domain,
	})

	return csrfHandler
}