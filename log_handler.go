package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

type loggingHandler struct {
	http.ResponseWriter
	statusCode int
	contentLen int
}

const (
	timeFormat = "02/Jan/2006:15:04:05 -0700"
)

//LoggingHandler returns a http.Handler that wraps next, and prints requests and responses in Apache Combined Log Format.
func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &loggingHandler{w, http.StatusOK, 0}
		next.ServeHTTP(lw, r)
		printLog(lw.statusCode, lw.contentLen, r)
	})
}

func printLog(status, length int, r *http.Request) {
	fmt.Fprintf(output, "%s %s %s %s %s %s %s %d %d %s %s\n", r.RemoteAddr, "-", "-", time.Now().Format(timeFormat), r.Method, r.RequestURI, r.Proto, status, length, r.Referer(), r.UserAgent())
}

//WriteHeader shadows http.ResponseWriter.WriteHeader()
func (l *loggingHandler) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

//Write shadows http.ResponseWriter.Write
func (l *loggingHandler) Write(b []byte) (n int, err error) {
	n, err = l.ResponseWriter.Write(b)
	l.contentLen += n
	return
}
