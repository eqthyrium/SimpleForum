package customHttp

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
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

	files := []string{"../ui/html/homepage.tmpl.html"}
	handler.homePageGet(w, r, files)

}

func (handler *HandlerHttp) homePageGet(w http.ResponseWriter, r *http.Request, files []string) {

	var userId int
	role := r.Context().Value("Role").(string)
	if role != "Guest" {
		userId = r.Context().Value("UserId").(int)
	} else {
		userId = -1
	}
	requestedCategories := r.URL.Query()["categories"]
	myposts := r.URL.Query().Get("myposts")
	mylikeposts := r.URL.Query().Get("mylikedposts")
	requestedmoderation := r.FormValue("requestmoderation")
	//report := r.FormValue("report")

	categories, err := handler.Service.GetAllCategories()
	if err != nil {

		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of getting all categories", err))
		serverError(w)
		return

	}

	if role == "User" && requestedmoderation == "true" {
		err := handler.Service.RequestToBeModerator(userId)
		if errors.Is(err, domain.ErrRepeatedRequest) {
			files = append(files, "../ui/html/error/requestmoderation.tmpl.html")
		} else if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of requesting moderator", err))
			serverError(w)
			return
		}

	}

	var posts []entity.Posts

	if role != "Guest" && (myposts == "true" || mylikeposts == "true") {

		if myposts == "true" {
			posts, err = handler.Service.GetMyCreatedPosts(userId)
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of getting my created posts", err))
				serverError(w)
				return
			}
		} else {
			posts, err = handler.Service.GetMyLikedPosts(userId)
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of getting my liked posts", err))
				serverError(w)
				return
			}
		}

	} else {
		posts, err = handler.Service.GetLatestPosts(requestedCategories)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of getting latest posts", err))
			serverError(w)
			return
		}
	}

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

	//fmt.Println("Incoming userId:", userId, "\nIncoming role:", role)
	//fmt.Println("Inserting csrfText:", csrfText)
	//fmt.Println("And its map:", session.MapUUID[userId])

	CSRFMap[session.MapUUID[userId]] = csrfText

	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "homepage", map[string]interface{}{
		"CSRFText":           csrfText,
		"Categories":         categories,
		"Posts":              posts,
		"UserIdentification": userId,
		"Role":               role,
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
