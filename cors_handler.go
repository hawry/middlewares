package middlewares

import (
	"net/http"
	"strings"
)

var (
	allowedCORSMethods = []string{}
	allowedOrigins     = []string{}
	allowedHeaders     = []string{}
	exposedHeaders     = []string{}
)

var (
	corsAccessControlRequestMethod    = "Access-Control-Request-Method"  //to let the server know which method will be used in the upcoming request
	corsAccessControlRequestHeaders   = "Access-Control-Request-Headers" //used when issuing a preflight req to let the server know what HTTP headers will be used in the upcoming request
	corsOrigin                        = "Origin"                         //origin of the request or preflight request, only server name, can be an empty string (e.g. when the source is a data URL)
	corsVary                          = "Vary"
	corsAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	corsAccessControlExposeHeaders    = "Access-Control-Expose-Headers"    // lets a server whitelist headers that browsers are allowed to access
	corsAccessControlMaxAge           = "Access-Control-Max-Age"           // how long the results of a preflight request <delta-seconds>
	corsAccessControlAllowCredentials = "Access-Control-Allow-Credentials" // when used as part of preflight response, indicates if the actual request can be made using credentials
	corsAccessControlAllowMethods     = "Access-Control-Allow-Methods"     // method or methods allowed for the resource <method>[, <method>]*
	corsAccessControlAllowHeaders     = "Access-Control-Allow-Headers"     // indicate which http headers that can be used when making the actual request
)

var (
	supportCredentials bool
)

func init() {
	supportCredentials = false //disallow credentials by default
}

//AllowCORSMethods specified which methods that are allowed for CORS
func AllowCORSMethods(methods ...string) {
	allowedCORSMethods = append(allowedCORSMethods, methods...)
}

//AllowCORSOrigins specifies which origins to allow
func AllowCORSOrigins(origins ...string) {
	allowedOrigins = append(allowedOrigins, origins...)
}

//AllowCORSHeaders specifies which headers that can be used
func AllowCORSHeaders(headers ...string) {
	allowedHeaders = append(allowedHeaders, headers...)
}

//AllowCORSExposedHeaders whitelists which headers are allowed for the browsers to access
func AllowCORSExposedHeaders(headers ...string) {
	exposedHeaders = append(exposedHeaders, headers...)
}

//SupportCredentials sets the Support-Credentials header to b
func SupportCredentials(b bool) {
	supportCredentials = b
}

//CORSHandler appends CORS headers to the response if any CORS headers are present in the request or as a preflight request
func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get(corsOrigin)
		if !isAllowedOrigin(origin) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set(corsAccessControlAllowOrigin, origin) // we allow this origin

		if r.Method == http.MethodOptions {
			if isAllowedMethod(r.Header.Get(corsAccessControlRequestMethod)) {
				w.Header().Set(corsAccessControlAllowMethods, strings.Join(allowedCORSMethods, ", "))
			}
			w.Header().Set(corsAccessControlMaxAge, "0")

			//Request-headers are only set on preflight requests
			sv := getAllowedHeaders(strings.Split(r.Header.Get(corsAccessControlRequestHeaders), ","))
			if len(sv) > 0 {
				w.Header().Set(corsAccessControlAllowHeaders, strings.Join(sv, ", "))
			}

			if len(exposedHeaders) > 0 {
				w.Header().Set(corsAccessControlExposeHeaders, strings.Join(exposedHeaders, ", "))
			}
		}

		if supportCredentials {
			w.Header().Set(corsAccessControlAllowCredentials, "true")
		} else {
			if allowAll() {
				w.Header().Set(corsVary, "Origin") // tell the client the origin might be changing depending on who's asking
			}
		}
		next.ServeHTTP(w, r)
	})
}

func getAllowedHeaders(req []string) (all []string) {
	all = []string{}
	for _, s := range req {
		for _, c := range allowedHeaders {
			if strings.ToLower(strings.TrimSpace(s)) == strings.ToLower(strings.TrimSpace(c)) {
				all = append(all, c)
			}
		}
	}
	return
}

func allowAll() bool {
	for _, ao := range allowedOrigins {
		if ao == "*" {
			return true
		}
	}
	return false
}

func isAllowedOrigin(o string) bool {
	if o == "" {
		return false
	}
	for _, ao := range allowedOrigins {
		if ao == o || ao == "*" {
			return true
		}
	}
	return false
}

func isAllowedMethod(m string) bool {
	if m == "" {
		return false
	}
	for _, am := range allowedCORSMethods {
		if am == m || am == "*" {
			return true
		}
	}
	return false
}
