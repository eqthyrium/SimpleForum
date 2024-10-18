package customHttp

import (
	"SimpleForum/internal/service/session"
	"SimpleForum/pkg/logger"
	"fmt"
	"html/template"
	"net/http"
)

func (handler *HandlerHttp) homePage(w http.ResponseWriter, r *http.Request) {
	handler.DebugLog.Println("homePage handler is activated")

	if r.URL.Path != "/" {
		clientError(w, nil, http.StatusNotFound, nil)
		handler.InfoLog.Println("incorrect request's endpoint")
		return
	}
	if r.Method == http.MethodGet {
		role := r.Context().Value("Role").(string)

		var userId int
		if role != "Guest" {
			userId = r.Context().Value("UserId").(int)
		}
		if role == "Admin" {
			fmt.Fprintf(w, "<html><body><h1>Welcome to the Admin Page!!!!!!!!!!!!!</h1></body></html>")

		} else if role == "Moderator" {
			fmt.Fprintf(w, "<html><body><h1>Welcome to the Moderator Page!!!!!!!!!</h1></body></html>")

		} else if role == "User" {

			files := []string{
				"../ui/html/homepage.html",
			}

			tmpl, err := template.ParseFiles(files...)
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePage", "There is a problem in the process of parsing the html files with template", err))
				serverError(w)
				return
			}

			csrfText, err := session.GenerateRandomCSRFText()
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePage", "There is a problem in the process of generating random CSRF text", err))
				serverError(w)
				return
			}

			CSRFMap[session.MapUUID[userId]] = csrfText

			err = tmpl.Execute(w, map[string]interface{}{
				"CSRFText": csrfText,
			})

			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePage", "There is a problem in the process of executing the template", err))
				serverError(w)
				return
			}

		} else if role == "Guest" {
			fmt.Fprintf(w, "<html><body><h1>Welcome to the Guest Page!!!!!</h1></body></html>")
		}
	}

}
