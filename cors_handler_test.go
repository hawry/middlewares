package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllowAll(t *testing.T) {
	t.Run("Return true on allow all", func(t *testing.T) {
		allowedOrigins = []string{"*"}
		assert.True(t, allowAll(), "should be true")
	})

	t.Run("Return false on allow all", func(t *testing.T) {
		allowedOrigins = []string{"http://example.com"}
		assert.False(t, allowAll(), "should be false")
	})
}

func TestAllowedOrigin(t *testing.T) {
	t.Run("Return true if origin exists", func(t *testing.T) {
		allowedOrigins = []string{"http://www.hawry.net", "http://www.benefactory.se"}
		assert.True(t, isAllowedOrigin("http://www.benefactory.se"), "should be true")
	})

	t.Run("Return false if origin doesn't exists", func(t *testing.T) {
		allowedOrigins = []string{"http://www.hawry.net", "http://www.benefactory.se"}
		assert.False(t, isAllowedOrigin("http://www.example.com"), "should be false")
		assert.False(t, isAllowedOrigin(""), "should be false")
	})
}

func TestAllowedMethod(t *testing.T) {
	t.Run("Return false if method isn't specified", func(t *testing.T) {
		allowedCORSMethods = []string{"GET", "POST"}
		assert.False(t, isAllowedMethod("DELETE"), "should be false")
		assert.False(t, isAllowedMethod(""), "should be false")
	})

	t.Run("Return true if method is specified", func(t *testing.T) {
		allowedCORSMethods = []string{"GET", "POST"}
		assert.True(t, isAllowedMethod("GET"), "should be true")
	})
}

func TestAllowedHeaders(t *testing.T) {
	allowedHeaders = []string{"X-Real-IP", "Content-Type"}
	reqHds := []string{"X-Real-IP", "X-Requested-With"}
	shouldBe := []string{"X-Real-IP"}
	assert.Subset(t, getAllowedHeaders(reqHds), shouldBe)
}

func TestSetMethods(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", "/index", nil)
	if err != nil {
		t.Fatal(err)
	}
	//Clear any leftover data from other tests
	allowedCORSMethods = []string{}
	allowedHeaders = []string{}
	allowedOrigins = []string{}
	exposedHeaders = []string{}
	AllowCORSOrigins("http://example.com", "http://localhost")
	AllowCORSMethods("GET", "POST")
	AllowCORSHeaders("X-Real-IP", "Content-Type")
	AllowCORSExposedHeaders("X-Exposed-Header")
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	t.Run("Requesting allowed origin", func(t *testing.T) {
		req.Header.Set(corsOrigin, "http://localhost")
		rr := httptest.NewRecorder()
		handler := CORSHandler(corsHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, "http://localhost", rr.Header().Get(corsAccessControlAllowOrigin), "should be equal")
	})

	t.Run("Requesting disallowed origin", func(t *testing.T) {
		req.Header.Set(corsOrigin, "http://notexamples")
		rr := httptest.NewRecorder()
		handler := CORSHandler(corsHandler)
		handler.ServeHTTP(rr, req)
		assert.Empty(t, rr.Header().Get(corsAccessControlAllowOrigin), "should be empty")
	})

	t.Run("Preflight request with methods and headers", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodOptions, "/index", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set(corsOrigin, "http://localhost")
		req.Header.Set(corsAccessControlRequestMethod, "GET")
		req.Header.Set(corsAccessControlRequestHeaders, "X-Real-IP")
		rr := httptest.NewRecorder()
		handler := CORSHandler(corsHandler)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, "http://localhost", rr.Header().Get(corsAccessControlAllowOrigin))
		assert.Equal(t, "GET, POST", rr.Header().Get(corsAccessControlAllowMethods))
		assert.Equal(t, "X-Exposed-Header", rr.Header().Get(corsAccessControlExposeHeaders))
		assert.Equal(t, "X-Real-IP", rr.Header().Get(corsAccessControlAllowHeaders))
	})

}

func TestSupportCredentials(t *testing.T) {
	allowedCORSMethods = []string{}
	allowedHeaders = []string{}
	allowedOrigins = []string{}
	exposedHeaders = []string{}

	req, err := http.NewRequest("GET", "/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Origin", "http://localhost")

	SupportCredentials(true)
	AllowCORSOrigins("*")

	rr := httptest.NewRecorder()
	handler := CORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, "true", rr.Header().Get(corsAccessControlAllowCredentials))
	assert.NotEqual(t, "*", rr.Header().Get(corsAccessControlAllowOrigin))
	// assert.Equal(t, "Origin", rr.Header().Get(corsVary))
	assert.Equal(t, "http://localhost", rr.Header().Get(corsAccessControlAllowOrigin))
}

func ExampleCORSHandler() {

	AllowCORSOrigins("http://example.com")  //requests from http://example.com are allowed
	AllowCORSMethods("GET", "POST", "HEAD") //requests with method GET, POST and HEAD are allowed
	AllowCORSHeaders("X-Requested-With")    //the server will be able to handle the header X-Requested-With

	corsHandler := CORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ... do something
	}))

	http.Handle("/", corsHandler)
	http.ListenAndServe(":8080", nil)
}
