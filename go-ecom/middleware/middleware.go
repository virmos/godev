package middleware

import (
	"go-ecom/data"

	"github.com/tsawler/celeritas"
)

type Middleware struct {
	App *celeritas.Celeritas
	Models data.Models
}