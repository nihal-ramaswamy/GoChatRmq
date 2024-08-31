package dto

import (
	"net/http"
)

type Routes struct {
	Path    string
	Handler http.Handler
}
