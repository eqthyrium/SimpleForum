package customHttp

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
)

func (handler *HandlerHttp) createPost(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("homePage handler is activated")

	if r.URL.Path != "/create/post" {
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
		files := []string{"../ui/html/createpostpage.tmpl.html"}
		handler.createPostPage(w, r, files)
	}

	role := r.Context().Value("Role").(string)

	if role == "Guest" {
		customLogger.DebugLogger.Println("The guest is trying to creating a post")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		userId := r.Context().Value("UserId").(int)
		title := r.FormValue("title")
		content := r.FormValue("content")
		requestedCategories := r.URL.Query()["categories"]
		fmt.Println("content:", content)
		err := handler.Service.CreatePost(userId, title, content, requestedCategories)
		if err != nil {
			customLogger.ErrorLogger.Println(err)
			serverError(w)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func (handler *HandlerHttp) createPostPage(w http.ResponseWriter, r *http.Request, files []string) {
	userId := r.Context().Value("UserId").(int)

	categories, err := handler.Service.GetAllCategories()
	if err != nil {

		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "createPostPage", "There is a problem in the process of getting all categories", err))
		serverError(w)
		return

	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "createPostPage", "There is a problem in the process of parsing the html files with template", err))
		serverError(w)
		return
	}

	csrfText, err := session.GenerateRandomCSRFText()
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "createPostPage", "There is a problem in the process of generating random CSRF text", err))
		serverError(w)
		return
	}

	CSRFMap[session.MapUUID[userId]] = csrfText

	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "createpostpage", map[string]interface{}{
		"CSRFText":   csrfText,
		"Categories": categories,
	})
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "createPostPage", "There is a problem in the process of rendering template to the buffer", err))
		serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = buf.WriteTo(w)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "createPostPage", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
		serverError(w)
		return
	}
}
