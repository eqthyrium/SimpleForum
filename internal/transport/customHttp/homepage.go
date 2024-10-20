package customHttp

import (
	"SimpleForum/internal/service/session"
	"SimpleForum/pkg/logger"
	"bytes"
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
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		handler.InfoLog.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	if r.Method == http.MethodGet {
		var userId int
		role := r.Context().Value("Role").(string)
		if role != "Guest" {
			userId = r.Context().Value("UserId").(int)
		}
		switch role {
		case "Admin":
			homePageGet(w, userId, []string{"../ui/html/homePage/homepage.html", "../ui/html/homePage/homePageAdmin.tmpl.html"})
		case "Moderator":
			homePageGet(w, userId, []string{"../ui/html/homePage/homepage.html", "../ui/html/homePage/homePageModerator.tmpl.html"})
		case "User":
			homePageGet(w, userId, []string{"../ui/html/homePage/homepage.html", "../ui/html/homePage/homePageUser.tmpl.html"})
		case "Guest":

			fmt.Fprintf(w, "<html><body><h1>Welcome to the Guest Page!!!!!</h1></body></html>")
		}
	}

	if r.Method == http.MethodPost {
		//
	}

}

func homePageGet(w http.ResponseWriter, userId int, files []string) {

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of parsing the html files with template", err))
		serverError(w)
		return
	}

	csrfText, err := session.GenerateRandomCSRFText()
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of generating random CSRF text", err))
		serverError(w)
		return
	}

	CSRFMap[session.MapUUID[userId]] = csrfText
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "homepage", map[string]interface{}{
		"CSRFText": csrfText,
	})
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of rendering template to the buffer", err))
		serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = buf.WriteTo(w)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
		serverError(w)
		return
	}
}
