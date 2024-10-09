package middleware

import (
	"net/http"
)

func LogInRoleDispatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value("Role").(string)

		switch role {
		case "Guest":
			// LogIn webpage
		default:
			// UnAuthorized error webpage
		}
	})
}
