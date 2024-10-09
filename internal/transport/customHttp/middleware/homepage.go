package middleware

import (
	"net/http"
)

func RoleDispatcherMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
