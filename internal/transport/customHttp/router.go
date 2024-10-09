package customHttp

import "net/http"

func (handler *HandlerHttp) Routering() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.homePage)
	mux.HandleFunc("/login", handler.logIn)
	mux.HandleFunc("/signup", handler.signUp)

	return http.HandlerFunc(mux.ServeHTTP)
}
