package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractBearer(t *testing.T) {
	req, err := http.NewRequest("GET", "/tokens", nil)
	if err != nil {
		t.Fatal(err)
	}
	raw := "thisisabearertokenthatshouldbefound"
	req.Header.Add("Authentication", fmt.Sprintf("Bearer %s", raw))
	tval, err := bearer(req.Header)
	assert.Nil(t, err, "should be nil")
	assert.Equal(t, raw, tval, "should be equal")

	t.Run("Empty bearer", func(t *testing.T) {
		req.Header.Set("Authentication", "Bearer ")
		tval, err := bearer(req.Header)
		assert.NotNil(t, err, "should not be nil")
		assert.Empty(t, tval, "should be empty")
		assert.Equal(t, "empty bearer token", err.Error(), "should be equal")
	})

	t.Run("Empty authentication", func(t *testing.T) {
		req.Header.Del("Authentication")
		tval, err := bearer(req.Header)
		assert.NotNil(t, err, "should not be nil")
		assert.Empty(t, tval, "should be empty")
		assert.Equal(t, "no authentication header found", err.Error(), "should be equal")
	})
}

func TestTokenHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/tokens", nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Empty token", func(t *testing.T) {
		req.Header.Set("Authentication", "Bearer")

		ctxHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := Token(r.Context())
			assert.NotNil(t, err, "should not be nil")
			assert.Empty(t, token, "should be empty")
			assert.Equal(t, "no token found in context", err.Error(), "should be equal")
		})

		rr := httptest.NewRecorder()
		handler := TokenHandler(ctxHandler)
		handler.ServeHTTP(rr, req)
	})

	t.Run("With token", func(t *testing.T) {
		req.Header.Set("Authentication", "Bearer thisisatoken")

		ctxHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := Token(r.Context())
			assert.Nil(t, err, "should be nil")
			assert.Equal(t, "thisisatoken", token, "should be equal")
		})

		rr := httptest.NewRecorder()
		handler := TokenHandler(ctxHandler)
		handler.ServeHTTP(rr, req)
	})
}

func ExampleToken() {
	defaultHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := Token(r.Context())
		if err != nil {
			// error handling
		}
		fmt.Printf("%s", token)
	})
	// ...
	http.Handle("/", defaultHandler)
	// Output: 1/mZ1edKKACtPAb7zGlwSzvs72PvhAbGmB8K1ZrGxpcNM (example)
}
func ExampleTokenHandler() {
	defaultHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do something
	})

	http.Handle("/", TokenHandler(defaultHandler))
	http.ListenAndServe(":3000", nil)
}
