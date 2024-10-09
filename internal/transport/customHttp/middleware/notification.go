package middleware

import (
	"net/http"
)

func NotificationRoleDispatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value("Role").(string)

		switch role {
		case "Admin":
			//Admin notification webpage handler
		case "Moderator":
			// Moderator notification webpage handler
		case "User":
			// User notification webpage handler
		default:

			// UnAuthorized error
		}
	})
}
