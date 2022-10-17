package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// mux.Use(NoSurf)

	mux.Post("/api/authenticate", app.CreateAuthToken)
	mux.Post("/api/is-authenticated", app.CheckAuthentication)

	mux.Route("/pusher", func(mux chi.Router) {
		mux.Use(app.Auth)
		mux.Post("/auth", app.PusherAuth)
	})

	mux.Route("/api/admin", func(mux chi.Router) {
		mux.Use(app.Auth)

		//  future use: message when user deleted
		mux.Get("/private-message", app.SendPrivateMessage)
		
		// excel
		mux.Post("/upload-excel", app.PostExcel)
		mux.Post("/download-excel", app.GetExcel)

		// settings
		mux.Post("/settings", app.PostSettings)

		// users
		mux.Post("/user/{id}", app.PostOneUser)
		mux.Post("/user/delete/{id}", app.DeleteUser)
		
		// preferences
		mux.Post("/preference/ajax/set-system-pref", app.SetSystemPref)
		mux.Post("/preference/ajax/toggle-monitoring", app.ToggleMonitoring)

		// hosts
		mux.Post("/host/{id}", app.PostHost)
		mux.Post("/host/delete/{id}", app.DeleteHost)
		mux.Post("/host/ajax/toggle-service", app.ToggleServiceForHost)
		mux.Post("/perform-check/{id}/{oldStatus}", app.TestCheck)

		// uptime report
		mux.Post("/send-range-uptime-report", app.SendRangeUptimeReport)
		mux.Post("/send-range-uptime-report-cached", app.SendRangeUptimeReportCached)
	})

	return mux
}
