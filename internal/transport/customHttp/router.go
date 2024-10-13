package customHttp

import (
	"SimpleForum/internal/transport/customHttp/middleware"
	"net/http"
)

func (handler *HandlerHttp) Routering() http.Handler {

	defaultPostMiddleware := func(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
		return middleware.LoggingMiddleware(middleware.SecurityMiddleware(middleware.CSRFPostMiddleware(middleware.RoleAdjusterMiddleware(http.HandlerFunc(next)))))
	}
	defaultMiddleware := func(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
		return middleware.LoggingMiddleware(middleware.SecurityMiddleware(middleware.RoleAdjusterMiddleware(http.HandlerFunc(next))))
	}

	mux := http.NewServeMux()
	mux.Handle("/", defaultMiddleware(handler.homePage))
	mux.Handle("/login", defaultPostMiddleware(handler.logIn))
	mux.Handle("/signup", defaultPostMiddleware(handler.signUp))
	//mux.Handle("/post", postPagePath)
	return http.HandlerFunc(mux.ServeHTTP)
}
