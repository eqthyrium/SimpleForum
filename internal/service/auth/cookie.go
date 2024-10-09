package auth

import (
	"net/http"
)

func SetTokenToCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Secure:   false, // Set to true in production with HTTPS
	}
	http.SetCookie(w, cookie)
}

func GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
