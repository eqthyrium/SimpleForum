package customHttp

import (
	"net/http"
)

func (handler *HandlerHttp) Routering() http.Handler {
	Middleware := func(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
		return LoggingMiddleware(SecurityMiddleware(PanicMiddleware(RoleAdjusterMiddleware(CSRFMiddleware((http.HandlerFunc(next)))))))
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/auth/login", Middleware(handler.logIn))
	mux.Handle("/auth/signup", Middleware(handler.signUp))
	mux.Handle("/logout", Middleware(handler.logOut))
	mux.Handle("/oauth2/google", Middleware(handler.googleAuthentication))
	mux.Handle("/oauth2/google/callback", Middleware(handler.googleCallback))
	mux.Handle("/oauth2/github", Middleware(handler.githubAuthentication))
	mux.Handle("/oauth2/github/callback", Middleware(handler.githubCallback))
	mux.Handle("/", Middleware(handler.homePage))
	mux.Handle("/reaction", Middleware(handler.reaction))
	mux.Handle("/post/", Middleware(handler.postPage))
	mux.Handle("/create/post", Middleware(handler.createPost))
	mux.Handle("/activity", Middleware(handler.myActivityPage))

	return http.HandlerFunc(mux.ServeHTTP)
}
