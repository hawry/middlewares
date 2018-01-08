package middlewares

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type basicAuth struct {
	user string
	pass string
}

//BasicAuthorizationHandler extracts any Authorization info of type Basic.
func BasicAuthorizationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hash, err := basicAuthHash(r.Header)
		if err != nil {
			fmt.Fprintf(output, "warn: %s\n", err.Error())
			next.ServeHTTP(w, r)
			return
		}
		plain, err := decodeBase64(hash)
		if err != nil {
			fmt.Fprintf(output, "warn: %s\n", err.Error())
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), authContextKey, newBasicAuth(plain))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//BasicCredentials returns the user, pass or an error if the values could not be handled or doesn't exist in the context
func BasicCredentials(ctx context.Context) (user, pass string, err error) {
	auth, ok := ctx.Value(authContextKey).(basicAuth)
	if !ok {
		return "", "", errors.New("no basic authorization data could be found")
	}
	return auth.user, auth.pass, nil
}

func basicAuthHash(h http.Header) (b64 string, err error) {
	auth := h.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		return "", errors.New("no basic authorization header found")
	}
	hash := auth[6:]
	if !(len(hash) > 0) {
		return "", errors.New("empty basic authorization hash")
	}
	return hash, nil
}

func decodeBase64(b64 string) (plain string, err error) {
	d, err := base64.URLEncoding.DecodeString(b64)
	return string(d), err
}

func newBasicAuth(plain string) (a basicAuth) {
	a = basicAuth{}
	s := strings.Split(plain, ":")
	if len(s) != 2 {
		return
	}
	if !(len(s[0]) > 0) || !(len(s[1]) > 0) {
		return
	}
	a.user = s[0]
	a.pass = s[1]
	return
}
