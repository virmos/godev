package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/virmos/django/mailer"
	"github.com/virmos/django"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes
	a.use(a.App.CheckForMaintenanceMode)

	// add routes here
	a.get("/", a.Handlers.Home)
	a.get("/go-page", a.Handlers.GoPage)
	a.get("/jet-page", a.Handlers.JetPage)
	a.get("/sessions", a.Handlers.SessionTest)

	// user api
	a.post("/api/register", a.Handlers.RegisterUser)
	a.post("/api/login", a.Handlers.LoginUser)
	a.post("/api/update-user", a.Handlers.UpdateUser)
	a.post("/api/delete-user-by-id", a.Handlers.DeleteUser)
	a.post("/api/get-users", a.Handlers.GetAllUsers)
	a.post("/api/get-user-by-id", a.Handlers.GetUserById)
	a.post("/api/get-user-by-email", a.Handlers.GetUserByEmail)

	a.get("/form", a.Handlers.Form)
	a.post("/api/form", a.Handlers.PostForm)

	// filesystem api
	a.get("/upload", a.Handlers.DjangoUpload)
	a.post("/upload", a.Handlers.PostDjangoUpload)

	a.get("/list-fs", a.Handlers.ListFS)
	a.get("/delete-from-fs", a.Handlers.DeleteFromFS)
	
	a.get("/json", a.Handlers.JSON)
	a.get("/xml", a.Handlers.XML)
	a.get("/download-file", a.Handlers.DownloadFile)

	a.get("/crypto", a.Handlers.TestCrypto)

	// cache api
	a.get("/cache-test", a.Handlers.ShowCachePage)
	a.post("/api/save-in-cache", a.Handlers.SaveInCache)
	a.post("/api/get-from-cache", a.Handlers.GetFromCache)
	a.post("/api/delete-from-cache", a.Handlers.DeleteFromCache)
	a.post("/api/empty-cache", a.Handlers.EmptyCache)

	a.get("/test-mail", func(w http.ResponseWriter, r *http.Request) {
		msg := mailer.Message{
			From:        "reportspeech@gmail.com",
			To:          "hadonggiang1810@gmail.com",
			Subject:     "Test Subject - sent using an API",
			Template:    "test",
			Attachments: nil,
			Data:        nil,
		}

		a.App.Mail.Jobs <- msg
		res := <-a.App.Mail.Results
		if res.Error != nil {
			a.App.ErrorLog.Println(res.Error)
		}
		// err := a.App.Mail.SendSMTPMessage(msg)
		// if err != nil {
		// 	a.App.ErrorLog.Println(err)
		// 	return
		// }

		fmt.Fprint(w, "Sent mail!")
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	// routes from django
	a.App.Routes.Mount("/django", django.Routes())
	// a.App.Routes.Mount("/api", a.ApiRoutes())

	return a.App.Routes
}
