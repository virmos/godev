package middleware

import (
	"ecom/data"

	"github.com/virmos/django"
)

type Middleware struct {
	App *django.Django
	Models data.Models
}