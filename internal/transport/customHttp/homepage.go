package customHttp

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
	"bytes"
	"errors"
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

		files := []string{"../ui/html/homepage.html"}

		requestedCategories := r.URL.Query()["categories"]

		if role != "Guest" {
			myposts := r.URL.Query().Get("myposts")
			mylikeposts := r.URL.Query().Get("mylikedposts")

			if myposts == "true" {

			}

			if mylikeposts == "true" {

			}

		}

		// Here you have also take from the paramaters mypost, myliked argument

		categories, err := handler.Service.GetAllCategories()
		if err != nil {
			if errors.Is(err, domain.ErrCategoryNotFound) {
				// I have to output error message to the client in order to notify that there is no any kind of categories
			} else {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePage", "There is a problem in the process of getting all categories", err))
				serverError(w)
				return
			}
		}

		posts, err := handler.Service.GetLatestPosts(requestedCategories)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePage", "There is a problem in the process of getting latest posts", err))
			serverError(w)
			return
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePage", "There is a problem in the process of parsing the html files with template", err))
			serverError(w)
			return
		}

		var buf bytes.Buffer

		if role != "Guest" {

			csrfText, err := session.GenerateRandomCSRFText()
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePage", "There is a problem in the process of generating random CSRF text", err))
				serverError(w)
				return
			}

			CSRFMap[session.MapUUID[userId]] = csrfText

			err = tmpl.ExecuteTemplate(&buf, "homepage", map[string]interface{}{
				"CSRFText":           csrfText,
				"Categories":         categories,
				"Posts":              posts,
				"UserIdentification": userId,
				"Role":               role,
			})

		} else {
			// There is a call for usecase for getting the data about categories, and a list of all posts

			err = tmpl.ExecuteTemplate(&buf, "homepage", map[string]interface{}{
				"Categories":         categories,
				"Posts":              posts,
				"UserIdentification": -1,
				"Role":               role,
			})

		}

		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePage", "There is a problem in the process of rendering template to the buffer", err))
			serverError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = buf.WriteTo(w)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePage", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
			serverError(w)
			return
		}

	}

	if r.Method == http.MethodPost {
		//
	}

}
