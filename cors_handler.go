package middlewares

var (
	possibleCORSMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	allowedCORSMethods  = []string{}
)

//AllowCORSMethods specified which methods that are allowed for CORS
func AllowCORSMethods(methods ...string) {
	allowedCORSMethods = append(allowedCORSMethods, methods...)
}
