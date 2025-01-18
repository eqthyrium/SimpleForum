package customHttp

import (
	"fmt"
	"net/http"
)

func (handler *HandlerHttp) results(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	categories := r.Form["categories"] // Get selected categories
	fmt.Println("These are the list of incoming categories:", categories)
}
