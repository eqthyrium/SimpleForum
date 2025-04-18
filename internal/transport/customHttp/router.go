package customHttp

import (
	"net/http"
	"time"
)

func (handler *HandlerHttp) Routering() http.Handler {
	rateLimiter := RateLimiterMiddleware(5, 3*time.Second)

	Middleware := func(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
		return rateLimiter(LoggingMiddleware(SecurityMiddleware(PanicMiddleware(RoleAdjusterMiddleware(CSRFMiddleware((http.HandlerFunc(next))))))))
	}
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
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
	mux.Handle("/notification", Middleware(handler.notification))
	mux.Handle("/editing", Middleware(handler.editing))
	mux.Handle("/myactivity", Middleware(handler.myActivityPage))
	mux.Handle("/categorylist", Middleware(handler.categoryListPage))
	mux.Handle("/moderationlist", Middleware(handler.moderationList))

	return http.HandlerFunc(mux.ServeHTTP)
}
