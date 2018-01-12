package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tNotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("the page could not be found"))
})

var tFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("the page could be found"))
})

func TestErrorBody(t *testing.T) {
	req, err := http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := ErrorHandler(tNotFound)
	handler.ServeHTTP(rr, req)
	assert.Contains(t, rr.Body.String(), "the page could not be found")
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPassthrough(t *testing.T) {
	req, err := http.NewRequest("GET", "/found", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := ErrorHandler(tFound)
	handler.ServeHTTP(rr, req)

	assert.Contains(t, rr.Body.String(), "the page could be found")
	assert.Equal(t, http.StatusOK, rr.Code)
}
