package customHttp

import (
	"bytes"
	"html/template"
	"net/http"

	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
)

func (handler *HandlerHttp) myActivityPage(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("myActivity page is activated")
	if r.URL.Path != "/myactivity" {
		clientError(w, nil, http.StatusNotFound, nil)
		customLogger.InfoLogger.Println("incorrect request's endpoint, it's requested endpoint is:", r.URL.Path)
		return
	}
	if r.Method != http.MethodGet {
		customLogger.InfoLogger.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	role := r.Context().Value("Role").(string)

	if role == "Guest" {
		customLogger.DebugLogger.Println("The guest is trying to creating a post")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userId := r.Context().Value("UserId").(int)

	// requestedCategories := r.URL.Query()["categories"]

	getAllMyPosts, err := handler.Service.GetMyCreatedPosts(userId)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "myActivityPage", "There is a problem in the process of getting all categories", err))
		serverError(w)
		return
	}

	getAllMyLikedPosts, err := handler.Service.GetMyLikedPosts(userId)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "myActivityPage", "There is a problem in the process of getting all categories", err))
		serverError(w)
		return
	}

	getDislikedPost, err := handler.Service.GetMyDislikedPosts(userId)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "myActivityPage", "There is a problem in the process of getting all categories", err))
		serverError(w)
		return
	}

	// for _, r := range getDislikedPost{
	// 	r.DislikeCount
	// }

	getCommentedPosts, err := handler.Service.GetMyCommentedPosts(userId)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "myActivityPage", "There is a problem in the process of getting all categories", err))
		serverError(w)
		return
	}

	// var postId int
	// for _, r := range getCommentedPosts {
	// 	postId = r.PostId
	// }

	getComments, err := handler.Service.GetComments(userId)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "myActivityPage", "There is a problem in the process of getting all categories", err))
		serverError(w)
		return
	}
	type Post struct {
		PostId  int
		Title   string
		Content string
		Comment string
		Image   string
	}
	var postsWithComments []Post

	for _, post := range getCommentedPosts {
		for _, comment := range getComments {
			if post.PostId == comment.PostId {
				postsWithComments = append(postsWithComments, Post{
					PostId:  post.PostId,
					Title:   post.Title,
					Content: post.Content,
					Comment: comment.Content,
					Image:   post.Image,
				})
			}
		}
	}

	files := []string{"../ui/html/myactivity.html"}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "myActivityPage", "There is a problem in the process of parsing the html files with template", err))
		serverError(w)
		return
	}

	csrfText, err := session.GenerateRandomCSRFText()
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "myActivityPage", "There is a problem in the process of generating random CSRF text", err))
		serverError(w)
		return
	}

	CSRFMap[session.MapUUID[userId]] = csrfText

	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "myactivity.html", map[string]interface{}{
		"PostContent":      getAllMyPosts,
		"likedPosts":       getAllMyLikedPosts,
		"dislikedPosts":    getDislikedPost,
		"myCommentedPosts": postsWithComments,
	})
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "myActivityPage", "There is a problem in the process of rendering template to the buffer", err))
		serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = buf.WriteTo(w)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "myActivityPage", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
		serverError(w)
		return
	}
}
