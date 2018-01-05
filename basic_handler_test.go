package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashExtractor(t *testing.T) {
	req, err := http.NewRequest("GET", "/basic", nil)
	if err != nil {
		t.Fatal(err)
	}
	raw := "dG9tOnNoYXJkd2FyZQ==" //plain: tom:shardware
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", raw))

	rval, err := basicAuthHash(req.Header)
	assert.Nil(t, err, "should be nil")
	assert.Equal(t, raw, rval, "should be equal")

	t.Run("Empty authorization hash", func(t *testing.T) {
		req.Header.Set("Authorization", "Basic ")
		rval, err := basicAuthHash(req.Header)
		assert.NotNil(t, err, "should not be nil")
		assert.Empty(t, rval, "should be empty")
		assert.Equal(t, "empty basic authorization hash", err.Error(), "should be equal")
	})

	t.Run("No authorization header", func(t *testing.T) {
		req.Header.Del("Authorization")
		rval, err := basicAuthHash(req.Header)
		assert.NotNil(t, err, "should not be nil")
		assert.Empty(t, rval, "should be empty")
		assert.Equal(t, "no basic authorization header found", err.Error(), "should be equal")
	})

	t.Run("Decode b64", func(t *testing.T) {
		rval, err := decodeBase64(raw)
		assert.Nil(t, err, "should be nil")
		assert.Equal(t, "tom:shardware", rval, "should be equal")

		_, err = decodeBase64("notbase64")
		assert.NotNil(t, err)
	})

	t.Run("BasicAuth struct", func(t *testing.T) {
		plain := "tom:shardware"
		a := newBasicAuth(plain)
		assert.Equal(t, "tom", a.user, "should be equal")
		assert.Equal(t, "shardware", a.pass, "should be equal")

		plain = "tomshardware"
		a = newBasicAuth(plain)
		assert.Empty(t, a.user)
		assert.Empty(t, a.pass)

		plain = "tom:"
		a = newBasicAuth(plain)
		assert.Empty(t, a.user)
		assert.Empty(t, a.pass)

		plain = ":shardware"
		a = newBasicAuth(plain)
		assert.Empty(t, a.user)
		assert.Empty(t, a.pass)
	})
}

func TestBasicAuthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/basic", nil)
	if err != nil {
		t.Fatal(err)
	}
	b64 := "dG9tOnNoYXJkd2FyZQ==" //plain: tom:shardware

	t.Run("Extracts from context", func(t *testing.T) {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b64))
		ctxHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, err := BasicCredentials(r.Context())
			assert.Nil(t, err, "should be nil")
			assert.Equal(t, "tom", u)
			assert.Equal(t, "shardware", p)
		})
		rr := httptest.NewRecorder()
		h := BasicAuthenticationHandler(ctxHandler)
		h.ServeHTTP(rr, req)
	})

	t.Run("Missing basic auth data", func(t *testing.T) {
		req.Header.Del("Authorization")
		ctxHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, err := BasicCredentials(r.Context())
			assert.NotNil(t, err, "should not be nil")
			assert.Empty(t, u)
			assert.Empty(t, p)
		})

		var b bytes.Buffer
		SetOutput(&b)

		rr := httptest.NewRecorder()
		h := BasicAuthenticationHandler(ctxHandler)
		h.ServeHTTP(rr, req)

		bs := string(b.Bytes())
		assert.Equal(t, "warn: no basic authorization header found", bs)
	})

	t.Run("Invalid basic auth data", func(t *testing.T) {
		req.Header.Add("Authorization", "Basic thisisnotbase64")
		ctxHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, err := BasicCredentials(r.Context())
			assert.NotNil(t, err, "should not be nil")
			assert.Empty(t, u)
			assert.Empty(t, p)
		})

		var b bytes.Buffer
		SetOutput(&b)

		rr := httptest.NewRecorder()
		h := BasicAuthenticationHandler(ctxHandler)
		h.ServeHTTP(rr, req)

		bs := string(b.Bytes())
		assert.Contains(t, bs, "warn: ")
	})
}

func ExampleBasicCredentials() {
	defaultHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, err := BasicCredentials(r.Context())
		if err != nil {
			// error handling
		}
		fmt.Printf("%s, %s", user, pass)
	})
	// ...
	http.Handle("/", defaultHandler)
}

func ExampleBasicAuthenticationHandler() {
	defaultHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do something
	})

	http.Handle("/", BasicAuthenticationHandler(defaultHandler))
	http.ListenAndServe(":3000", nil)
}
