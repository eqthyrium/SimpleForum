package customHttp

import (
	"SimpleForum/internal/domain"
	"errors"
	"fmt"
	"net/http"
)

func (handler *Handler) logIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		// Think about error handling, and logging it properly
		handler.notFound(w)
		return
	}
	if !(r.Method != http.MethodPost || r.Method != http.MethodGet) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodGet {

	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			handler.serverError(w, err)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		err = handler.Service.LogIn(email, password)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				// Here must be webpage of error input for LogIning.
			} else {
				handler.serverError(w, fmt.Errorf("Http-logIn: %w", err))
				return
			}
		}

		// Think about session based on token here
	}

}
