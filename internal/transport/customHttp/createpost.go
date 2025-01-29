package customHttp

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	"SimpleForum/internal/config"
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
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
	role := r.Context().Value("Role").(string)

	if role == "Guest" {
		customLogger.DebugLogger.Println("The guest is trying to creating a post")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		files := []string{"./ui/html/createpostpage.tmpl.html"}
		handler.createPostPage(w, r, files)
	}

	if r.Method == http.MethodPost {
		userId := r.Context().Value("UserId").(int)
		title := r.FormValue("title")
		content := r.FormValue("content")
		requestedCategories := r.Form["categories"]

		err := r.ParseMultipartForm(config.MaxImageSize)
		if err != nil {
			customLogger.InfoLogger.Println(logger.ErrorWrapper("UseCase", "createPost", "Failed to parse the multipart form", err))
			clientError(w, nil, http.StatusBadRequest, nil)
			return
		}

		var myfile *entity.MyFile
		// Retrieve the file
		file, header, err := r.FormFile("image")

		if errors.Is(err, http.ErrMissingFile) { // No image uploaded
			myfile = nil
		} else if err != nil {
			customLogger.InfoLogger.Println(logger.ErrorWrapper("UseCase", "createPost", "Failed to retrieve the image file", err))
			clientError(w, nil, http.StatusBadRequest, nil)
			return
		} else {
			myfile = &entity.MyFile{
				FileContent: file,
				FileHeader:  header,
			}
		}

		err = handler.Service.CreatePost(userId, title, content, requestedCategories, myfile)
		files := []string{"./ui/html/createpostpage.tmpl.html"}

		if errors.Is(err, domain.ErrNoCategories) {
			files = append(files, "./ui/html/error/nocategories.tmp.html")
			handler.createPostPage(w, r, files)
			return
		} else if errors.Is(err, domain.ErrNotValidContent) {
			files = append(files, "./ui/html/error/postpagecontent.tmpl.html")
			handler.createPostPage(w, r, files)
			return
		} else if errors.Is(err, domain.ErrLargeImageSize) {
			files = append(files, "./ui/html/error/imagesize.tmpl.html")
			handler.createPostPage(w, r, files)
			return
		} else if errors.Is(err, domain.ErrInvalidImageType) {
			files = append(files, "./ui/html/error/imagetype.tmpl.html")
			handler.createPostPage(w, r, files)
			return
		} else if err != nil {
			customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "createPost", "Failed to create a post", err))
			serverError(w)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		err = file.Close()
		if err != nil {
			customLogger.ErrorLogger.Println(logger.ErrorWrapper("Transport", "createPost", "Failed to close the file", err))
			serverError(w)
			return
		}

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
