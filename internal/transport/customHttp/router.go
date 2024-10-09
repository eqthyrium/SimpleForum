package customHttp

import (
	"SimpleForum/internal/transport/customHttp/middleware"
	"net/http"
)

func (handler *HandlerHttp) Routering() http.Handler {
	homePagePath := middleware.LoggingMiddleware(middleware.SecurityMiddleware(middleware.RoleAdjusterMiddleware(middleware.HomePageRoleDispatcher(nil))))
	logInPath := middleware.LoggingMiddleware(middleware.SecurityMiddleware(middleware.RoleAdjusterMiddleware(middleware.LogInRoleDispatcher(nil))))
	signUpPath := middleware.LoggingMiddleware(middleware.SecurityMiddleware(middleware.RoleAdjusterMiddleware(middleware.SignUpRoleDispatcher(nil))))
	postPagePath := middleware.LoggingMiddleware(middleware.SecurityMiddleware(middleware.RoleAdjusterMiddleware(middleware.PostRoleDispatcher(nil))))

	mux := http.NewServeMux()
	mux.Handle("/", homePagePath)
	mux.Handle("/post", postPagePath)
	mux.Handle("/login", logInPath)
	mux.Handle("/signup", signUpPath)

	return http.HandlerFunc(mux.ServeHTTP)
}
