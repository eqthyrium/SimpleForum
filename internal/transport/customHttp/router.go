package customHttp

import (
	"net/http"
)

func (handler *HandlerHttp) Routering() http.Handler {

	defaultCSRFMiddleware := func(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
		return LoggingMiddleware(SecurityMiddleware(RoleAdjusterMiddleware((http.HandlerFunc(next)))))
	}
	defaultMiddleware := func(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
		return LoggingMiddleware(SecurityMiddleware(RoleAdjusterMiddleware(http.HandlerFunc(next))))
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/", defaultMiddleware(handler.homePage))
	mux.Handle("/auth/login", defaultCSRFMiddleware(handler.logIn))
	mux.Handle("/auth/signup", defaultCSRFMiddleware(handler.signUp))
	mux.Handle("/logout", defaultCSRFMiddleware(handler.logOut))
	// Example of serving static files

	//mux.Handle("/post", postPagePath)
	return http.HandlerFunc(mux.ServeHTTP)
}
