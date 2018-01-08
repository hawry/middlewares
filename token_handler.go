package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

//TokenHandler extracts any Bearer tokens from the request
func TokenHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := bearer(r.Header)
		if err != nil {
			fmt.Fprintf(output, "warn: %s\n", err.Error())
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), tokenContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//Token returns any Bearer tokens that was found using TokenHandler, if no token is found it returns an error
func Token(ctx context.Context) (token string, err error) {
	if t, ok := ctx.Value(tokenContextKey).(string); ok {
		return t, nil
	}
	return "", errors.New("no token found in context")
}

func bearer(h http.Header) (token string, err error) {
	auth := h.Get("Authentication")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", errors.New("no authentication header found")
	}
	token = auth[7:] // The string 'Bearer' and a whitespace "Bearer "
	if !(len(token) > 0) {
		return "", errors.New("empty bearer token")
	}
	return token, nil
}
