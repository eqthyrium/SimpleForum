package customHttp

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"SimpleForum/internal/domain"
	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
)

func (handler *HandlerHttp) editing(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("editing handler is activated")

	if r.URL.Path != "/editing" {
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
		customLogger.DebugLogger.Println("The guest is trying to edit")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		files := []string{"./ui/html/editingpage.tmpl.html"}
		handler.editingGet(w, r, files)
		return
	}
	postId := r.FormValue("postId")
	commentId := r.FormValue("commentId")
	content := r.FormValue("content")
	userId := r.Context().Value("UserId").(int)

	if r.Method == http.MethodPost {
		if !(commentId != "" && postId != "") && (commentId != "" || postId != "") { // commentId != ""XOR postId != ""
			if postId != "" {
				postNumber, err := strconv.Atoi(postId)
				if err != nil {
					customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editing", "Failed to convert the value from string to int type", err))
					serverError(w)
					return
				}
				if postNumber < 0 || postNumber > 4e9 {
					customLogger.InfoLogger.Println("The incoming request's body data is not correct")
					clientError(w, nil, http.StatusBadRequest, nil)
					return
				}
				err = handler.Service.EditCertainPost(userId, postNumber, content)
				if errors.Is(err, domain.ErrPostNotFound) {
					customLogger.DebugLogger.Println("The client is trying to edit the post which is not related to him")
					http.Redirect(w, r, "/", http.StatusSeeOther)

				} else if errors.Is(err, domain.ErrNotValidContent) {
					files := []string{"./ui/html/editingpage.tmpl.html"}
					files = append(files, "./ui/html/error/postpagecontent.tmpl.html")
					handler.editingGet(w, r, files)

				} else if err != nil {
					customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editing", "Failed to check the correspondence of the user to a post", err))
					serverError(w)

				}

			} else {
				commentNumber, err := strconv.Atoi(commentId)
				if err != nil {
					customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editing", "Failed to convert the value from string to int type", err))
					serverError(w)
					return
				}
				if commentNumber < 0 || commentNumber > 4e9 {
					customLogger.InfoLogger.Println("The incoming request's body data is not correct")
					clientError(w, nil, http.StatusBadRequest, nil)
					return
				}
				err = handler.Service.EditCertainCommentary(userId, commentNumber, content)
				if errors.Is(err, domain.ErrCommentaryNotFound) {
					customLogger.DebugLogger.Println("The client is trying to edit the comment which is not related to him")
					http.Redirect(w, r, "/", http.StatusSeeOther)

				} else if errors.Is(err, domain.ErrNotValidContent) {
					files := []string{"./ui/html/editingpage.tmpl.html"}
					files = append(files, "./ui/html/error/postpagecontent.tmpl.html")
					handler.editingGet(w, r, files)

				} else if err != nil {
					customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editing", "Failed to check the correspondence of the user to a comment", err))
					serverError(w)
				}

			}

			http.Redirect(w, r, "/", http.StatusSeeOther) // Fallback to homepage

		} else {
			customLogger.InfoLogger.Println("The incoming request's body data is not correct")
			clientError(w, nil, http.StatusBadRequest, nil)
			return
		}
	}
}

func (handler *HandlerHttp) editingGet(w http.ResponseWriter, r *http.Request, files []string) {
	userId := r.Context().Value("UserId").(int)

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

	postId := r.FormValue("postId")
	commentId := r.FormValue("commentId")

	if !(commentId != "" && postId != "") && (commentId != "" || postId != "") { // commentId != ""XOR postId != ""

		if postId != "" {
			postNumber, err := strconv.Atoi(postId)
			if err != nil {
				customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editingGet", "Failed to convert the value from string to int type", err))
				serverError(w)
				return
			}
			if postNumber < 0 || postNumber > 4e9 {
				customLogger.InfoLogger.Println("The incoming request's body data is not correct")
				clientError(w, nil, http.StatusBadRequest, nil)
				return
			}
			err = handler.Service.EditCertainPost(userId, postNumber, "")
			if errors.Is(err, domain.ErrPostNotFound) {
				customLogger.DebugLogger.Println("The client is trying to edit the post which is not related to him")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else if !errors.Is(err, domain.ErrNotValidContent) && err != nil {
				customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editingGet", "Failed to check the correspondence of the user to a post", err))
				serverError(w)
				return
			}

			postInfo, err := handler.Service.GetCertainPostInfo(postNumber)
			if err != nil {
				customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editingGet", "Failed to retrieve a post info", err))
				serverError(w)
				return
			}

			err = tmpl.ExecuteTemplate(&buf, "editingpage", map[string]interface{}{
				"CSRFText": csrfText,
				"PostInfo": postInfo,
			})
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "editingGet", "There is a problem in the process of rendering template to the buffer", err))
				serverError(w)
				return
			}

		} else {
			commentNumber, err := strconv.Atoi(commentId)
			if err != nil {
				customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editingGet", "Failed to convert the value from string to int type", err))
				serverError(w)
				return
			}
			if commentNumber < 0 || commentNumber > 4e9 {
				customLogger.InfoLogger.Println("The incoming request's body data is not correct")
				clientError(w, nil, http.StatusBadRequest, nil)
				return
			}
			err = handler.Service.EditCertainCommentary(userId, commentNumber, "")
			if errors.Is(err, domain.ErrCommentaryNotFound) {
				customLogger.DebugLogger.Println("The client is trying to edit the comment which is not related to him")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else if !errors.Is(err, domain.ErrNotValidContent) && err != nil {
				customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editingGet", "Failed to check the correspondence of the user to a commentary", err))
				serverError(w)
				return
			}

			commentInfo, err := handler.Service.GetCertainCommentaryInfo(commentNumber)
			if err != nil {
				customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "editingGet", "Failed to retrieve a commentary's info", err))
				serverError(w)
				return
			}

			err = tmpl.ExecuteTemplate(&buf, "editingpage", map[string]interface{}{
				"CSRFText":    csrfText,
				"CommentInfo": commentInfo,
			})
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "editingGet", "There is a problem in the process of rendering template to the buffer", err))
				serverError(w)
				return
			}

		}
	} else {
		customLogger.InfoLogger.Println("Incoming request's body data is not correct to the logic of the application")
		clientError(w, nil, http.StatusBadRequest, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = buf.WriteTo(w)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "editingGet", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
		serverError(w)
		return
	}
}
