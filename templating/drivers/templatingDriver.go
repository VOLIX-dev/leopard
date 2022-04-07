package drivers

import (
	"github.com/gorilla/mux"
	"io"
)

type TemplatingDriver interface {
	// RenderTemplate renders a template with the given data.
	RenderTemplate(template string, writer io.Writer, data map[string]Value) error

	Load(path string, router *mux.Router) error
}

type Value interface{}
