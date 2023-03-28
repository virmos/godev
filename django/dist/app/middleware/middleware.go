package middleware

import (
	"app/data"

	"github.com/virmos/django"
)

type Middleware struct {
	App *django.Django
	Models data.Models
}