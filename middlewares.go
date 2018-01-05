// Package middlewares aims to create a set of commonly used middleware http.Handlers for use with the default http package. All handlers only takes a http.Handler as an argument, and returns only http.Handler, to more easily be chained with handler chain libraries (e.g. https://github.com/justinas/alice)
package middlewares

type contextKey string

const (
	jwtContextKey  contextKey = "mw_jwt_context_key"
	fpContextKey   contextKey = "mw_fp_context_key"
	authContextKey contextKey = "mw_auth_context_key"
)
