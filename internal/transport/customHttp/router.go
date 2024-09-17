package customHttp

import "net/http"

func (handler *Handler) Routering() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.homePage)
	mux.HandleFunc("/login", handler.login)
	mux.HandleFunc("/signup", handler.signUp)

	return http.HandlerFunc(mux.ServeHTTP)
}
