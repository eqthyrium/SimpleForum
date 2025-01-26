package customHttp

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"SimpleForum/internal/domain"
	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
)

func (handler *HandlerHttp) postPage(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("postPage handler is activated")

	path := r.URL.Path
	if !strings.HasPrefix(path, "/post/") {
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		customLogger.InfoLogger.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	postID := strings.TrimPrefix(path, "/post/")
	if postID == "" {
		http.Error(w, "Post ID not provided", http.StatusBadRequest)
		return
	}

	number, err := strconv.Atoi(postID)
	if err != nil {
		customLogger.InfoLogger.Println(err)
		clientError(w, nil, http.StatusBadRequest, err)
		return
	}
	if number <= 0 || number > 4e9 {
		customLogger.InfoLogger.Println("The incoming request's post Id number is not corresponds to application logic")
		clientError(w, nil, http.StatusBadRequest, nil)
		return
	}

	if r.Method == http.MethodGet {
		files := []string{"../ui/html/postpage.tmpl.html"}
		handler.postPageGet(w, r, files, number)
	}

	// For like/dislike operation
	if r.Method == http.MethodPost {
		var userId int
		role := r.Context().Value("Role").(string)
		if role != "Guest" {
			userId = r.Context().Value("UserId").(int)
		} else {
			userId = -1
		}

		if role == "Guest" {
			customLogger.DebugLogger.Println("The guest is trying to comment to a post")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		}

		commentaried := r.FormValue("commentary")
		commentText := r.FormValue("commentText")
		if commentaried == "true" {
			err := handler.Service.CreateCommentary(userId, number, commentText)

			if errors.Is(err, domain.ErrNotValidContent) {
				var files []string
				files = append(files, "../ui/html/postpage.tmpl.html")
				files = append(files, "../ui/html/error/postpagecontent.tmpl.html")
				handler.postPageGet(w, r, files, number)
				return

			} else if err != nil {
				customLogger.ErrorLogger.Println(err)
				serverError(w)
				return
			}

		} else {
			customLogger.InfoLogger.Println("The request's commentary parameter is incorrect")
			clientError(w, nil, http.StatusBadRequest, nil)
			return
		}
		referer := r.Header.Get("Referer")
		if referer != "" {
			http.Redirect(w, r, referer, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther) // Fallback to homepage
		}
	}
}

func (handler *HandlerHttp) postPageGet(w http.ResponseWriter, r *http.Request, files []string, number int) {
	var userId int
	role := r.Context().Value("Role").(string)
	if role != "Guest" {
		userId = r.Context().Value("UserId").(int)
	} else {
		userId = -1
	}

	postNumberId := number

	// Think about not founding a particular post form the db
	postInfo, commentaries, err := handler.Service.GetCertainPostPage(postNumberId)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			customLogger.InfoLogger.Println(err)
			clientError(w, nil, http.StatusNotFound, err)
		} else {
			customLogger.ErrorLogger.Println(err)
			serverError(w)
		}
		return
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "postPageGet", "There is a problem in the process of parsing the html files with template", err))
		serverError(w)
		return
	}

	csrfText, err := session.GenerateRandomCSRFText()
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "postPageGet", "There is a problem in the process of generating random CSRF text", err))
		serverError(w)
		return
	}

	CSRFMap[session.MapUUID[userId]] = csrfText

	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "homepage", map[string]interface{}{
		"CSRFText":           csrfText,
		"PostContent":        postInfo,
		"Commentaries":       commentaries,
		"UserIdentification": userId,
		"Role":               role,
	})
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "postPageGet", "There is a problem in the process of rendering template to the buffer", err))
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
