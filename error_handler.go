package middlewares

import (
	"html/template"
	"log"
	"net/http"
)

type errorHandler struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

const errBody = `<!doctype HTML>
<html>
<head>
	<meta charset="utf-8"/>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.StatusCode}} - {{.StatusText}}</title>
	<style type="text/css">
	h1 {
		color:#666;
	}
	.content {
		text-align:center;
		margin-left: auto;
		margin-right: auto;
		max-width: 75%;
		font-size: 1.5rem;
	}
	.error-text {
		color:#666;
	}
	</style>
</head>
<body>

	<div class="content">
		<h1>{{.StatusCode}}</h1>
		<p class="error-text">
			{{.StatusMessage}}
		</p>
	</div>

</body>
</html>
`

var errTemplate *template.Template

func init() {
	errTemplate = template.Must(template.New("err").Parse(errBody))
}

//ErrorHandler will inject a html response to any error status code
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eh := &errorHandler{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(eh, r)
		if eh.statusCode >= 400 && eh.statusCode <= 509 {
			// do something
			dval := map[string]interface{}{
				"StatusCode":    eh.statusCode,
				"StatusText":    http.StatusText(eh.statusCode),
				"StatusMessage": string(eh.body),
			}
			errTemplate.Execute(w, dval)
		}
	})
}

// WriteHeader shadows ResponseWriter.Write
func (e *errorHandler) WriteHeader(code int) {
	e.statusCode = code
	e.ResponseWriter.WriteHeader(code)
}

//Write shadows http.ResponseWriter.Write
func (e *errorHandler) Write(b []byte) (n int, err error) {
	if e.statusCode >= 400 && e.statusCode <= 509 {
		log.Printf("err handler write called")
		e.body = append(e.body, b...) // use any written text as the error message
	} else {
		e.ResponseWriter.Write(b)
	}
	return
}
