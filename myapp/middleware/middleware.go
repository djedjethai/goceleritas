package middleware

import (
	"myapp/data"

	"github.com/djedjethai/celeritas"
)

type Middleware struct {
	App    *celeritas.Celeritas
	Models data.Models
}
