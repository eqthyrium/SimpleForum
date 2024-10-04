package customHttp

import (
	"SimpleForum/internal/domain"
	"errors"
	"fmt"
	"net/http"
)

func (handler *Handler) signUp(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {
		// Think about error handling, and logging it properly
		handler.notFound(w)
		return
	}
	if !(r.Method == http.MethodGet || r.Method == http.MethodPost) {
		// Think about error handling, and logging it properly
		handler.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodGet { // pass the singup webpage

	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			handler.serverError(w, err)
			return
		}

		nickname := r.FormValue("nickname")
		email := r.FormValue("email")
		password := r.FormValue("password")

		err = handler.Service.SignUp(nickname, email, password)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidCredential) {
				// Here we have to create an error webpage for a client
			} else {
				handler.serverError(w, fmt.Errorf("Http-signUP: %w", err))
				return
			}
		}

		// Here i have to give to the client the logIn website.
	}
}
