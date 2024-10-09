package middleware

import (
	"net/http"
)

func PostRoleDispatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value("Role").(string)

		switch role {
		case "Admin":
			//Admin post webpage handler
		case "Moderator":
			// Moderator post webpage handler
		case "User":
			// User post webpage handler
		default:
			// UnAuthorized error
		}
	})
}
