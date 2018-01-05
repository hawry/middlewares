package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "/index.html", nil)
	req.Header.Set("User-Agent", "MWTests")
	req.Header.Set("Referer", "testing")
	if err != nil {
		log.Fatal(err)
	}
	status := http.StatusOK
	conLen := 23
	var b bytes.Buffer
	SetOutput(&b)
	printLog(status, conLen, req)
	rval := string(b.Bytes())
	assert.Contains(t, rval, "GET  HTTP/1.1 200 23 testing MWTests")
}

func TestShadowResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	l := loggingHandler{rr, http.StatusOK, 0}
	l.WriteHeader(http.StatusBadGateway)
	sval := "this is a body"
	l.Write([]byte(sval))
	assert.Equal(t, http.StatusBadGateway, l.statusCode, "should be equal")
	assert.Equal(t, len(sval), l.contentLen, "should be equal")
}

func TestLoggingHandler(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("empty response")) //14 chars
	})

	notfoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	req, err := http.NewRequest("GET", "/index", nil)
	if err != nil {
		log.Fatal(err)
	}

	t.Run("200 OK Content Length", func(t *testing.T) {
		var b bytes.Buffer
		rr := httptest.NewRecorder()
		SetOutput(&b)
		handler := LoggingHandler(emptyHandler)
		handler.ServeHTTP(rr, req)
		rval := string(b.Bytes())
		assert.Contains(t, rval, "GET  HTTP/1.1 200 14")
	})

	t.Run("404 Not Found", func(t *testing.T) {
		var b bytes.Buffer
		rr := httptest.NewRecorder()
		SetOutput(&b)
		handler := LoggingHandler(notfoundHandler)
		handler.ServeHTTP(rr, req)
		rval := string(b.Bytes())
		assert.Contains(t, rval, "GET  HTTP/1.1 404 0")
	})

}

func ExampleLoggingHandler() {
	defaultHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do something
	})

	http.Handle("/", LoggingHandler(defaultHandler))
	http.ListenAndServe(":3000", nil)
}

func ExampleSetOutput() {
	//This will print the log to logFile
	logFile, _ := os.Open("logfile")
	SetOutput(logFile)
}

func ExampleSetOutput_multiWriter() {
	//This will print the log both to Stdout and a file. Any ínterface implementing io.Writer can be used.
	logFile, _ := os.Open("logfile")
	mw := io.MultiWriter(os.Stdout, logFile)
	SetOutput(mw)
}
