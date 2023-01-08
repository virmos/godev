package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	_ "cycir/internal/models"
	"net/http"
	"strings"
	"time"
)

// LoginScreen shows the home (login) screen
func (app *application) LoginScreen(w http.ResponseWriter, r *http.Request) {
	// if already logged in, take to dashboard
	if app.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
		return
	}
	err := app.RenderPage(w, r, "login", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Login attempts to log the user in
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	user, err := app.repo.GetUserByEmail(r.Form.Get("email"))
	if err != nil {
		app.errorLog.Println(err)
		ClientError(w, r, http.StatusBadRequest)
	}

	if r.Form.Get("remember") == "remember" {
		randomString := RandomString(12)
		hasher := sha256.New()

		_, err = hasher.Write([]byte(randomString))
		if err != nil {
			app.errorLog.Println(err)
			ClientError(w, r, http.StatusInternalServerError)
		}

		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		err = app.repo.InsertRememberMeToken(user.ID, sha)
		if err != nil {
			app.errorLog.Println(err)
			ClientError(w, r, http.StatusInternalServerError)
		}

		// write a cookie
		expire := time.Now().Add(365 * 24 * 60 * 60 * time.Second)
		cookie := http.Cookie{
			Name:     fmt.Sprintf("_%s_gowatcher_remember", app.PreferenceMap["identifier"]),
			Value:    fmt.Sprintf("%d|%s", user.ID, sha),
			Path:     "/",
			Expires:  expire,
			HttpOnly: true,
			Domain:   cfg.Domain,
			MaxAge:   315360000, // seconds in year
			Secure:   cfg.InProduction,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
	}

	// we authenticated. Get the user.
	u, err := app.repo.GetUserById(user.ID)
	if err != nil {
		app.errorLog.Println(err)
		ClientError(w, r, http.StatusInternalServerError)
		return
	}

	// get the token fetched from backend, every request, send this token
	token := r.Form.Get("token")
	expiry := r.Form.Get("expiry")

	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "hashedPassword", string(user.Password))
	app.Session.Put(r.Context(), "flash", "You've been logged in successfully!")
	app.Session.Put(r.Context(), "user", u)
	app.Session.Put(r.Context(), "token", string(token))
	app.Session.Put(r.Context(), "expiry", expiry)

	if r.Form.Get("target") != "" {
		http.Redirect(w, r, r.Form.Get("target"), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
}

// Logout logs the user out
func (app *application) Logout(w http.ResponseWriter, r *http.Request) {

	// delete the remember me token, if any
	cookie, err := r.Cookie(fmt.Sprintf("_%s_gowatcher_remember", app.PreferenceMap["identifier"]))
	if err != nil {
	} else {
		key := cookie.Value
		// have a remember token, so get the token
		if len(key) > 0 {
			// key length > 0, so it might be a valid token
			split := strings.Split(key, "|")
			hash := split[1]
			err = app.repo.DeleteToken(hash)
			if err != nil {
				app.errorLog.Println(err)
				ClientError(w, r, http.StatusInternalServerError)
			}
		}
	}

	// delete the remember me cookie, if any
	delCookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_gowatcher_remember", app.PreferenceMap["identifier"]),
		Value:    "",
		Domain:   cfg.Domain,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	}
	http.SetCookie(w, &delCookie)

	_ = app.Session.RenewToken(r.Context())
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	app.Session.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
