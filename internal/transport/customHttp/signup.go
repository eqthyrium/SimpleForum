package customHttp

import "net/http"

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
			handler.clientError(w, http.StatusBadRequest)

			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
	}
}
