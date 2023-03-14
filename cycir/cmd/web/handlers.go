package main

import (
	_ "bytes"
	_ "cycir/internal/encryption"
	"cycir/internal/models"
	_ "cycir/internal/urlsigner"
	_ "encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/CloudyKit/jet/v6"
	"runtime/debug"
	"github.com/go-chi/chi/v5"
)

// AdminDashboard displays the dashboard
func (app *application) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	pending, healthy, warning, problem, err := app.repo.GetAllServiceStatusCounts()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("no_healthy", healthy)
	vars.Set("no_problem", problem)
	vars.Set("no_pending", pending)
	vars.Set("no_warning", warning)

	allHosts, err := app.repo.AllHosts()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	vars.Set("hosts", allHosts)

	err = app.RenderPage(w, r, "dashboard", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Events displays the events page
func (app *application) Events(w http.ResponseWriter, r *http.Request) {
	events, err := app.repo.GetAllEvents()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(jet.VarMap)
	data.Set("events", events)

	err = app.RenderPage(w, r, "events", data, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Settings displays the settings page
func (app *application) Settings(w http.ResponseWriter, r *http.Request) {
	err := app.RenderPage(w, r, "settings", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllHosts displays list of all hosts
func (app *application) AllHosts(w http.ResponseWriter, r *http.Request) {
	// get all hosts from database
	hosts, err := app.repo.AllHosts()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// send data to template
	vars := make(jet.VarMap)
	vars.Set("hosts", hosts)

	err = app.RenderPage(w, r, "hosts", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Host shows the host add/edit form
func (app *application) Host(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var h models.Host

	if id > 0 {
		// get the host from the database
		host, err := app.repo.GetHostByID(id)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		h = host
	}

	vars := make(jet.VarMap)
	vars.Set("host", h)

	err := app.RenderPage(w, r, "host", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllUsers lists all admin users
func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)

	u, err := app.repo.AllUsers()
	if err != nil {
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	vars.Set("users", u)

	err = app.RenderPage(w, r, "users", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// OneUser displays the add/edit user page
func (app *application) OneUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorLog.Println(err)
	}

	vars := make(jet.VarMap)

	if id > 0 {
		u, err := app.repo.GetUserById(id)
		if err != nil {
			ClientError(w, r, http.StatusBadRequest)
			return
		}
		vars.Set("user", u)
	} else {
		var u models.User
		vars.Set("user", u)
	}

	err = app.RenderPage(w, r, "user", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// ClientError will display error page for client error i.e. bad request
func ClientError(w http.ResponseWriter, r *http.Request, status int) {
	switch status {
	case http.StatusNotFound:
		show404(w, r)
	case http.StatusInternalServerError:
		show500(w, r)
	case http.StatusBadRequest:
		show400(w, r)
	default:
		http.Error(w, http.StatusText(status), status)
	}
}

// ServerError will display error page for internal server error
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(2, trace)
	show500(w, r)
}

func show400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/404.html")
}

func show404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/404.html")
}

func show500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/500.html")
}

func printTemplateError(w http.ResponseWriter, err error) {
	_, _ = fmt.Fprint(w, fmt.Sprintf(`<small><span class='text-danger'>Error executing template: %s</span></small>`, err))
}
