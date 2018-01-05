// Package middlewares aims to create a set of commonly used middleware http.Handlers for use with the default http package. All handlers only takes a http.Handler as an argument, and returns only http.Handler, to more easily be chained with handler chain libraries (e.g. https://github.com/justinas/alice)
package middlewares

import (
	"io"
	"os"
)

type contextKey string

var output io.Writer

const (
	tokenContextKey contextKey = "mw_token_context_key"
	fpContextKey    contextKey = "mw_fp_context_key"
	authContextKey  contextKey = "mw_auth_context_key"
)

func init() {
	output = os.Stdout
}

//SetOutput sets which io.Writer to print the log to. Default is os.Stdout.
//
// Setting an output will change the default output for ALL handlers in this package.
func SetOutput(w io.Writer) {
	output = w
}
