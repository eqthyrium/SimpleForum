package customHttp

import "net/http"

func (handler *Handler) login(w http.ResponseWriter, r *http.Request) {

	if !(r.Method != http.MethodPost || r.Method != http.MethodGet) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
