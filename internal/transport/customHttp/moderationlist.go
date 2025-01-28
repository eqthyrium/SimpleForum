package customHttp

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strconv"
)

func (handler *HandlerHttp) moderationList(w http.ResponseWriter, r *http.Request) {

	customLogger.DebugLogger.Println("moderationList handler is activated")

	if r.URL.Path != "/moderationlist" {
		clientError(w, nil, http.StatusNotFound, nil)
		customLogger.InfoLogger.Println("incorrect request's endpoint, it's requested endpoint is:", r.URL.Path)
		return
	}

	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		customLogger.InfoLogger.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	if role := r.Context().Value("Role").(string); role != "Admin" {
		customLogger.InfoLogger.Println("incorrect request's role")
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	files := []string{"../ui/html/moderatormanagement.tmpl.html"}
	if r.Method == http.MethodGet {
		userId := r.Context().Value("UserId").(int)

		users, err := handler.Service.GetCertainUsers()
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "moderationList", "There is an error of getting the certain users", err))
			serverError(w)
			return
		}
		reports, err := handler.Service.GetAllReports()
		if err != nil {
			customLogger.ErrorLogger.Print("Transport", "moderationList", "There is an error of getting all reports", err)
			serverError(w)
			return
		}
		requests, err := handler.Service.GetAllRequests()
		if err != nil {
			customLogger.ErrorLogger.Print("Transport", "moderationList", "There is an error of getting all requests", err)
			serverError(w)
			return
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "moderationList", "There is a problem in the process of parsing the html files with template", err))
			serverError(w)
			return
		}

		csrfText, err := session.GenerateRandomCSRFText()
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "moderationList", "There is a problem in the process of generating random CSRF text", err))
			serverError(w)
			return
		}

		CSRFMap[session.MapUUID[userId]] = csrfText

		var buf bytes.Buffer

		err = tmpl.ExecuteTemplate(&buf, "moderatormanagement", map[string]interface{}{
			"CSRFText": csrfText,
			"Users":    users,
			"Reports":  reports,
			"Requests": requests,
		})
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "moderationList", "There is a problem in the process of rendering template to the buffer", err))
			serverError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = buf.WriteTo(w)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "moderationList", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
			serverError(w)
			return
		}
	}

	if r.Method == http.MethodPost {
		promote := r.FormValue("promote")
		demote := r.FormValue("demote")
		deleting := r.FormValue("delete")
		accept := r.FormValue("accept")
		secondUserId := r.FormValue("userId")
		postId := r.FormValue("postId")
		secondUserIdNumber, err := strconv.Atoi(secondUserId)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "Failed to convert from string to integer value", err))
			serverError(w)
			return
		}
		if secondUserIdNumber < 0 || secondUserIdNumber > 4e9 {
			customLogger.InfoLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "Incoming request's userId is not corresponding to the application logic", err))
			clientError(w, nil, http.StatusBadRequest, nil)
			return
		}

		if promote == "true" {
			err := handler.Service.ChangeRole(secondUserIdNumber, "promote")
			if errors.Is(err, domain.ErrInvalidOperation) {
				customLogger.InfoLogger.Println(logger.ErrorWrapper("UseCase", "moderationList", "The incoming request's body is not correct, that's why it can not implement certain function", err))
				clientError(w, nil, http.StatusBadRequest, nil)
				return
			} else if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "There is a problem in the process of changing the role", err))
				serverError(w)
				return
			}
		} else if demote == "true" {
			err := handler.Service.ChangeRole(secondUserIdNumber, "demote")
			if errors.Is(err, domain.ErrInvalidOperation) {
				customLogger.InfoLogger.Println(logger.ErrorWrapper("UseCase", "moderationList", "The incoming request's body is not correct, that's why it can not implement certain function", err))
				clientError(w, nil, http.StatusBadRequest, nil)
				return
			} else if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "There is a problem in the process of changing the role", err))
				serverError(w)
				return
			}

		} else if deleting == "true" || deleting == "false" {
			postIdNumber, err := strconv.Atoi(postId)
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "Failed to convert from string to integer value", err))
				serverError(w)
				return
			}
			if postIdNumber < 0 || postIdNumber > 4e9 {
				customLogger.InfoLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "Incoming request's postId is not corresponding to the application logic", err))
				clientError(w, nil, http.StatusBadRequest, nil)
				return
			}

			if deleting == "true" {
				err := handler.Service.AcceptCertainReport(secondUserIdNumber, postIdNumber)
				if errors.Is(err, domain.ErrInvalidOperation) {
					customLogger.InfoLogger.Println(logger.ErrorWrapper("UseCase", "moderationList", "The incoming request's body is not correct, that's why it can not implement certain function", err))
					clientError(w, nil, http.StatusBadRequest, nil)
					return
				} else if err != nil {
					customLogger.ErrorLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "There is a problem in the process of accepting the report", err))
					serverError(w)
					return
				}
			} else {
				err := handler.Service.DeclineCertainReport(secondUserIdNumber, postIdNumber)
				if errors.Is(err, domain.ErrInvalidOperation) {
					customLogger.InfoLogger.Println(logger.ErrorWrapper("UseCase", "moderationList", "The incoming request's body is not correct, that's why it can not implement certain function", err))
					clientError(w, nil, http.StatusBadRequest, nil)
					return
				} else if err != nil {
					customLogger.ErrorLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "There is a problem in the process of declining the report", err))
					serverError(w)
					return
				}
			}

		} else if accept == "true" || accept == "false" {

			if accept == "true" {
				err := handler.Service.AcceptRequestToBeModerator(secondUserIdNumber)
				if errors.Is(err, domain.ErrInvalidOperation) {
					customLogger.InfoLogger.Println(logger.ErrorWrapper("UseCase", "moderationList", "The incoming request's body is not correct, that's why it can not implement certain function", err))
					clientError(w, nil, http.StatusBadRequest, nil)
					return
				} else if err != nil {
					customLogger.ErrorLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "There is a problem in the process of accepting the request", err))
					serverError(w)
					return
				}
			} else {
				err := handler.Service.DeclineRequestToBeModerator(secondUserIdNumber)
				if errors.Is(err, domain.ErrInvalidOperation) {
					customLogger.InfoLogger.Println(logger.ErrorWrapper("UseCase", "moderationList", "The incoming request's body is not correct, that's why it can not implement certain function", err))
					clientError(w, nil, http.StatusBadRequest, nil)
					return
				} else if err != nil {
					customLogger.ErrorLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "There is a problem in the process of declining the request", err))
					serverError(w)
					return
				}
			}

		} else {
			customLogger.InfoLogger.Print(logger.ErrorWrapper("UseCase", "moderationList", "Incoming request's userId is not corresponding to the application logic", err))
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
