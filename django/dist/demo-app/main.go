package main

import (
	"demo-app/data"
	"demo-app/handlers"
	"demo-app/middleware"

	"github.com/virmos/django"
)

type application struct {
	App *django.Django
	Handlers *handlers.Handlers
	Models data.Models
	Middleware *middleware.Middleware
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
