package customHttp

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/service/auth"
	"errors"
	"fmt"
	"net/http"
)

/*
How does CSRF process work
1) When the client gets the webpage where needs to be inserted data into it, the server must provide the webpage where its csrf token value and the embedded csrf unique string
must be same.
2) When the client sends (with Post method) data upon that fields on the webpage, the server must check whether the embedded data upon html is same to the data inside csrf token
If it is the server must provide the appropriate webpage resource to that request.
Otherwise the server must send back to the client error page 403 Forbidden

1) Generate CSRF unique number
2) I have to set it into the cookie as csrf token
3) Then it has to be embedded  into html, as hidden mark.
*/
func (handler *HandlerHttp) logIn(w http.ResponseWriter, r *http.Request) {
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
		// Here I have to send to the client's side against CSRF resolved  webpage

	}

	if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			handler.serverError(w, err)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		//flag := r.FormValue("flag")
		//
		//if flag {
		//	authentication()
		//}

		// Think about auth based on token here
		tokenSignature, _, err := handler.Service.LogIn(email, password)

		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				// Here must be webpage of error input for LogIning.
			} else {
				handler.serverError(w, fmt.Errorf("Http-logIn: %w", err))
				return
			}
		}

		// Cookies
		auth.SetTokenToCookie(w, tokenSignature)
		// Depending on the role, you have to return the appropriate webpage

	}

}
