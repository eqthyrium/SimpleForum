package customHttp

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

func (handler *HandlerHttp) homePage(w http.ResponseWriter, r *http.Request) {

	customLogger.DebugLogger.Println("homePage handler is activated")

	if r.URL.Path != "/" {
		clientError(w, nil, http.StatusNotFound, nil)
		customLogger.InfoLogger.Println("incorrect request's endpoint, it's requested endpoint is:", r.URL.Path)
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		customLogger.InfoLogger.Println("incorrect request's method")
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
			handler.homePageGet(w, r, userId, []string{"../ui/html/homePage/homepage.html", "../ui/html/homePage/homePageAdmin.tmpl.html"})
		case "Moderator":
			handler.homePageGet(w, r, userId, []string{"../ui/html/homePage/homepage.html", "../ui/html/homePage/homePageModerator.tmpl.html"})
		case "User":
			handler.homePageGet(w, r, userId, []string{"../ui/html/homePage/homepage.html", "../ui/html/homePage/homePageUser.tmpl.html"})
		case "Guest":
			handler.homePageGet(w, r, -1, []string{"../ui/html/homePage/homepage.html"})
			//fmt.Fprintf(w, "<html><body><h1>Welcome to the Guest Page!!!!!</h1></body></html>")
		}
	}

	if r.Method == http.MethodPost {
		//
	}

}

func (handler *HandlerHttp) homePageGet(w http.ResponseWriter, r *http.Request, userId int, files []string) {

	requestedCategories := r.URL.Query()["categories"]
	fmt.Println("The backend side, the incoming requestedCategories are: ", requestedCategories)
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of parsing the html files with template", err))
		serverError(w)
		return
	}

	var buf bytes.Buffer

	if userId >= 0 {

		csrfText, err := session.GenerateRandomCSRFText()
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of generating random CSRF text", err))
			serverError(w)
			return
		}

		CSRFMap[session.MapUUID[userId]] = csrfText

		err = tmpl.ExecuteTemplate(&buf, "homepage", map[string]interface{}{
			"CSRFText": csrfText,
		})
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of rendering template to the buffer", err))
			serverError(w)
			return
		}

	} else {
		// There is a call for usecase for getting the data about categories, and a list of all posts

		categories, err := handler.Service.GetAllCategories()
		if err != nil {
			if errors.Is(err, domain.ErrCategoryNotFound) {
				// I have to output error message to the client in order to notify that there is no any kind of categories
			} else {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of getting all categories", err))
				serverError(w)
				return
			}
		}
		fmt.Println("The categories from db, what i got:", categories)

		posts, err := handler.Service.GetLatestPosts(requestedCategories)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of getting latest posts", err))
			serverError(w)
			return
		}

		fmt.Println("The posts from db, what i got:", posts)
		err = tmpl.ExecuteTemplate(&buf, "homepage", map[string]interface{}{
			"Categories": categories,
			"Posts":      posts,
		})

		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of rendering template to the buffer", err))
			serverError(w)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, err = buf.WriteTo(w)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
		serverError(w)
		return
	}
}
