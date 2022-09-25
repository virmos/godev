package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	// default middleware
	mux.Use(SessionLoad)
	mux.Use(RecoverPanic)
	mux.Use(NoSurf)
	mux.Use(CheckRemember)

	// login
	mux.Get("/", app.LoginScreen)
	mux.Post("/", app.Login)

	mux.Get("/user/logout", app.Logout)

	mux.Route("/admin", func(mux chi.Router){
		// all admin routes are protected
		mux.Use(app.Auth)

		// overview
		mux.Get("/overview", app.AdminDashboard)

		// events
		mux.Get("/events", app.Events)

		// settings
		mux.Get("/settings", app.Settings)

		// service status pages (all hosts)
		mux.Get("/all-healthy", app.AllHealthyServices)
		mux.Get("/all-warning", app.AllWarningServices)
		mux.Get("/all-problems", app.AllProblemServices)
		mux.Get("/all-pending", app.AllPendingServices)

		// users
		mux.Get("/users", app.AllUsers)
		mux.Get("/user/{id}", app.OneUser)

		// hosts
		mux.Get("/host/all", app.AllHosts)
		mux.Get("/host/{id}", app.Host)
	})

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}