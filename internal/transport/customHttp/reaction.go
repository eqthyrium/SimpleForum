package customHttp

import (
	"SimpleForum/pkg/logger"
	"net/http"
	"strconv"
)

func (handler *HandlerHttp) reaction(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/reaction" {
		clientError(w, nil, http.StatusNotFound, nil)
		customLogger.InfoLogger.Println("incorrect request's endpoint, it's requested endpoint is:", r.URL.Path)
		return
	}

	if r.Method != http.MethodPost {
		customLogger.InfoLogger.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	var files []string
	role := r.Context().Value("Role").(string)
	if role == "Guest" {
		files = append(files, "../ui/html/homepage.tmpl.html")
		files = append(files, "../ui/html/error/homeroleforbidden.tmpl.html")
		handler.homePageGet(w, r, files)
		return
	}
	userId := r.Context().Value("UserId").(int)
	postId := r.FormValue("postId")
	like := r.FormValue("like")
	dislike := r.FormValue("dislike")
	commentId := r.FormValue("commentId")
	commented := r.FormValue("commented")

	if !(like == "true" && dislike == "true") && (like == "true" || dislike == "true") { // like == true XOR dislike == true

		if commented != "" {
			clientError(w, nil, http.StatusBadRequest, nil)
			return
		}

		var errr error

		if postId != "" {
			postIdInt, err := strconv.Atoi(postId)
			if err != nil {
				serverError(w)
				return
			}
			if like == "true" {
				errr = handler.Service.ExecutionOfReactionLD(userId, postIdInt, "post", "like")
			} else {
				errr = handler.Service.ExecutionOfReactionLD(userId, postIdInt, "post", "dislike")
			}
		} else if commentId != "" {
			commentIdInt, err := strconv.Atoi(commentId)
			if err != nil {
				serverError(w)
				return
			}

			if like == "true" {
				errr = handler.Service.ExecutionOfReactionLD(userId, commentIdInt, "comment", "like")
			} else {
				errr = handler.Service.ExecutionOfReactionLD(userId, commentIdInt, "comment", "dislike")
			}
		}

		if errr != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "reaction", "There is problem with implementing reaction", errr))
			serverError(w)
			return
		}

		referer := r.Header.Get("Referer")
		if referer != "" {
			http.Redirect(w, r, referer, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther) // Fallback to homepage
		}

	} else if commented != "" && (like == "" && dislike == "") {
		// Here i have to write the logic where the client writes comment for a particular post
	} else {
		// We have to output the error message that there is no corresponding body message from the client request
		clientError(w, nil, http.StatusBadRequest, nil)
		return
	}

}
