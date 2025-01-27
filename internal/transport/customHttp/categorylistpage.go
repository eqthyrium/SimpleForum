package customHttp

import (
	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
	"bytes"
	"html/template"
	"net/http"
	"strconv"
)

func (handler *HandlerHttp) categoryListPage(w http.ResponseWriter, r *http.Request) {

	customLogger.DebugLogger.Println("categoryListPage handler is activated")

	if r.URL.Path != "/categorylist" {
		customLogger.InfoLogger.Println("Incorrect request's endpoint:", r.URL.Path)
		clientError(w, nil, http.StatusNotFound, nil)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		customLogger.InfoLogger.Println("The request's method is not allowed:", r.URL.Path)
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	role := r.Context().Value("Role").(string)
	if role != "Admin" {
		customLogger.DebugLogger.Println("The guest is trying to access category management")
		clientError(w, nil, http.StatusForbidden, nil)
		return
	}
	// Обработка POST-запроса
	if r.Method == http.MethodPost {
		r.ParseForm()

		action := r.FormValue("action")
		categoryName := r.FormValue("categoryName")
		categoryId := r.FormValue("categoryId")

		if role == "Admin" {
			switch action {
			case "add":
				err := handler.Service.CreateCategory(categoryName)
				if err != nil {
					customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "categoryListPage", "Error creating category", err))
					serverError(w)
					return
				}
			case "delete":
				categoryIdNumber, err := strconv.Atoi(categoryId)
				if err != nil {
					customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "categoryListPage", "Error converting categoryId to int", err))
					return
				}
				if categoryIdNumber < 0 || categoryIdNumber > 4e9 {
					customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "categoryListPage", "Error converting categoryId is invalid", err))
					return
				}

				err = handler.Service.DeleteCategory(categoryIdNumber)
				if err != nil {
					customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "categoryListPage", "Error deleting category", err))
					serverError(w)
					return
				}
			default:
				clientError(w, nil, http.StatusBadRequest, nil)
				return
			}

			// После успешного выполнения действия перенаправляем обратно на страницу
			http.Redirect(w, r, "/categorylist", http.StatusSeeOther)
			return
		}

	}

	// Обработка GET-запроса
	if r.Method == http.MethodGet {

		userId := r.Context().Value("UserId").(int)

		categories, err := handler.Service.GetAllCategories()
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "categoryListPage", "Error getting all categories", err))
			serverError(w)
			return
		}

		files := []string{"../ui/html/categorylist.tmpl.html"}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "categoryListPage", "Error parsing the HTML template files", err))
			serverError(w)
			return
		}

		csrfText, err := session.GenerateRandomCSRFText()
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "categoryListPage", "Error generating CSRF token", err))
			serverError(w)
			return
		}

		CSRFMap[session.MapUUID[userId]] = csrfText

		var buf bytes.Buffer

		err = tmpl.ExecuteTemplate(&buf, "categorylist", map[string]interface{}{
			"CSRFText":           csrfText,
			"Categories":         categories,
			"UserIdentification": userId,
			"Role":               role,
		})
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "categoryListPage", "Error rendering template to the buffer", err))
			serverError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = buf.WriteTo(w)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "categoryListPage", "Error writing buffer data to the HTTP response", err))
			serverError(w)
			return
		}
	}

}
