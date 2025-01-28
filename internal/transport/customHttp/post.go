package customHttp

import (
	"SimpleForum/internal/config"
	"SimpleForum/internal/domain"
	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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

	postIdNumber, err := strconv.Atoi(postID)
	if err != nil {
		customLogger.InfoLogger.Println(err)
		clientError(w, nil, http.StatusBadRequest, err)
		return
	}
	if postIdNumber <= 0 || postIdNumber > 4e9 {
		customLogger.InfoLogger.Println("The incoming request's post Id postIdNumber is not corresponds to application logic")
		clientError(w, nil, http.StatusBadRequest, nil)
		return
	}

	if r.Method == http.MethodGet {
		files := []string{"../ui/html/postpage.tmpl.html"}
		handler.postPageGet(w, r, files, postIdNumber)
	}

	// For like/dislike operation
	if r.Method == http.MethodPost {
		role := r.Context().Value("Role").(string)

		if role == "Guest" {
			customLogger.DebugLogger.Println("The guest is trying to comment to a post")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		}
		userId := r.Context().Value("UserId").(int)

		commentaried := r.FormValue("commentary")
		commentText := r.FormValue("commentText")
		commentId := r.FormValue("commentId")
		deleting := r.FormValue("delete")
		report := r.FormValue("report")

		if !(commentaried == "true" && deleting == "true") && (commentaried == "true" || deleting == "true") { // like == true XOR dislike == true

			if commentaried == "true" {
				err := handler.Service.CreateCommentary(userId, postIdNumber, commentText)

				if errors.Is(err, domain.ErrNotValidContent) {
					var files []string
					files = append(files, "../ui/html/postpage.tmpl.html")
					files = append(files, "../ui/html/error/postpagecontent.tmpl.html")
					handler.postPageGet(w, r, files, postIdNumber)
					return

				} else if err != nil {
					customLogger.ErrorLogger.Println(err)
					serverError(w)
					return
				}

			} else if deleting == "true" {
				if commentId == "" {

					err := handler.Service.DeleteCertainPost(userId, postIdNumber, role)
					if err != nil {
						customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "postPage", "Failed to Delete a certain post", err))
						serverError(w)
						return
					}
				} else {
					commentIdNumber, err := strconv.Atoi(commentId)
					if err != nil {
						customLogger.ErrorLogger.Println(err)
						serverError(w)
						return
					}
					err = handler.Service.DeleteCertainCommentary(userId, commentIdNumber, role)
					if err != nil {
						customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "postPage", "Failed to Delete a certain commentary", err))
						serverError(w)
						return
					}
				}
			}
		} else {
			if role == "Moderator" && report == "true" {

				err := handler.Service.ReportPost(userId, postIdNumber)
				if errors.Is(err, domain.ErrRepeatedRequest) {
					referer := r.Header.Get("Referer")
					_, path, _ := strings.Cut(referer, "http://localhost"+*config.Config.Addr+"/")

					files := []string{"../ui/html/error/report.tmpl.html"}
					if path == "" {
						files = append(files, "../ui/html/homepage.tmpl.html")
						handler.homePageGet(w, r, files)
					} else {
						files = append(files, "../ui/html/postpage.tmpl.html")
						handler.postPageGet(w, r, files, postIdNumber)
					}

				} else if err != nil {
					customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of requesting moderator", err))
					serverError(w)
					return
				}
			} else {

				customLogger.InfoLogger.Println("The request's commentary parameter is incorrect")
				clientError(w, nil, http.StatusBadRequest, nil)
				return

			}

		}

		referer := r.Header.Get("Referer")
		if referer != "" && deleting != "true" {
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

	err = tmpl.ExecuteTemplate(&buf, "postpage", map[string]interface{}{
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
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "postPageGet", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
		serverError(w)
		return
	}
}
