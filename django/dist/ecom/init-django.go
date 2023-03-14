package main

import (
	"log"
	"ecom/data"
	"ecom/handlers"
	"ecom/middleware"
	"os"

	"github.com/virmos/django"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init django
	dj := &django.Django{}
	err = dj.New(path)
	if err != nil {
		log.Fatal(err)
	}

	dj.AppName = "ecom"

	myMiddleware := &middleware.Middleware{
		App: dj,
	}

	myHandlers := &handlers.Handlers{
		App: dj,
	}

	app := &application{
		App: dj,
		Handlers: myHandlers,
		Middleware: myMiddleware,
	}

	app.App.Routes = app.routes()

	app.Models = data.New(app.App.DB.Pool)
	myHandlers.Models = app.Models
	app.Middleware.Models = app.Models

	return app
}
