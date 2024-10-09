package customHttp

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (handler *HandlerHttp) serverError(w http.ResponseWriter, err error) {

	ErrorMessage := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	handler.ErrorLog.Print(ErrorMessage)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (handler *HandlerHttp) clientError(w http.ResponseWriter, statusCode int) {
	handler.ErrorLog.Print(http.StatusText(statusCode))
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (handler *HandlerHttp) notFound(w http.ResponseWriter) {
	handler.clientError(w, http.StatusNotFound)
}
