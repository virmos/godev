package main

import (
	"cycir/internal/models"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/justinas/nosurf"
	"log"
	"net/http"
	"time"
)

// TemplateData defines template data
type TemplateData struct {
	CSRFToken       string
	IsAuthenticated bool
	PreferenceMap   map[string]string
	User            models.User
	Token           string
	Flash           string
	Warning         string
	Error           string
	GwVersion       string
}

func addTemplateFunctions() {
	views.AddGlobal("humanDate", func(t time.Time) string {
		return HumanDate(t)
	})

	views.AddGlobal("dateFromLayout", func(t time.Time, l string) string {
		return FormatDateWithLayout(t, l)
	})

	views.AddGlobal("dateAfterYearOne", func(t time.Time) bool {
		return DateAfterY1(t)
	})
}

// HumanDate formats a time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// FormatDateWithLayout formats a time with provided (go compliant) format string, and returns it as a string
func FormatDateWithLayout(t time.Time, f string) string {
	return t.Format(f)
}

// DateAfterY1 is used to verify that a date is after the year 1 (since go hates nulls)
func DateAfterY1(t time.Time) bool {
	yearOne := time.Date(0001, 11, 17, 20, 34, 58, 651387237, time.UTC)
	return t.After(yearOne)
}

// views is the jet template set
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./cmd/web/views"),
	jet.InDevelopmentMode(),
)

// SetViews sets the view variable (for testing)
func SetViews(path string) {
	views = jet.NewSet(
		jet.NewOSFileSystemLoader(path),
		jet.InDevelopmentMode(),
	)
}

// DefaultData adds default data which is accessible to all templates
func (app *application) DefaultData(td TemplateData, r *http.Request) TemplateData {
	// refresh reference map every render
	preferenceMap = make(map[string]string)
	preferences, err := app.repo.AllPreferences()
	if err != nil {
		log.Fatal("Cannot read preferences:", err)
	}

	for _, pref := range preferences {
		app.PreferenceMap[pref.Name] = string(pref.Preference)
	}

	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = IsAuthenticated(r)
	td.PreferenceMap = app.PreferenceMap

	// if logged in, store user id in template data
	if td.IsAuthenticated {
		u := app.Session.Get(r.Context(), "user").(models.User)
		token := app.Session.Get(r.Context(), "token").(string)
		td.Token = token
		td.User = u
	}

	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	return td
}

// RenderPage renders a page using jet templates
func (app *application) RenderPage(w http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	// add default template data
	var td TemplateData
	if data != nil {
		td = data.(TemplateData)
	}

	// add default data
	td = app.DefaultData(td, r)

	// add template functions
	addTemplateFunctions()

	// load the template and render it
	t, err := views.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
