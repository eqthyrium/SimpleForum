package customHttp

import (
	"bytes"
	"html/template"
	"net/http"

	"SimpleForum/pkg/logger"
)

func (handler *HandlerHttp) notification(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("homePage handler is activated")

	if r.URL.Path != "/notification" {
		clientError(w, nil, http.StatusNotFound, nil)
		customLogger.InfoLogger.Println("incorrect request's endpoint, it's requested endpoint is:", r.URL.Path)
		return
	}

	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		customLogger.InfoLogger.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	role := r.Context().Value("Role").(string)
	if role == "Guest" {
		customLogger.DebugLogger.Println("The guest is trying to enter notification page")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userId := r.Context().Value("UserId").(int)

	if r.Method == http.MethodGet {
		notifications, err := handler.Service.GetNotifications(userId)
		if err != nil {
			customLogger.ErrorLogger.Println(err)
			serverError(w)
			return
		}

		files := []string{"./ui/html/notification.tmpl.html"}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "notification", "There is a problem in the process of parsing the html files with template", err))
			serverError(w)
			return
		}

		var buf bytes.Buffer

		err = tmpl.ExecuteTemplate(&buf, "notifications", map[string]interface{}{
			"Notifications": notifications,
		})
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "notification", "There is a problem in the process of rendering template to the buffer", err))
			serverError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = buf.WriteTo(w)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "notification", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
			serverError(w)
			return
		}
	}
}
