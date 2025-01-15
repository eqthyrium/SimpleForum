package customHttp

import (
	"SimpleForum/internal/config"
	"SimpleForum/internal/domain"
	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"strings"
)

var oauthState map[string]bool = make(map[string]bool)
var googleOauthConfig *oauth2.Config

func (handler *HandlerHttp) googleAuthentication(w http.ResponseWriter, r *http.Request) {

	customLogger.DebugLogger.Println("googleAuthentication handler is activated")
	googleOauthConfig = &oauth2.Config{
		ClientID:     config.Config.GoogleOauth.ClientID,
		ClientSecret: config.Config.GoogleOauth.ClientSecret,
		RedirectURL:  config.Config.GoogleOauth.RedirectURI,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	if r.URL.Path != "/oauth2/google" {
		customLogger.InfoLogger.Println("incorrect request's endpoint")
		clientError(w, nil, http.StatusNotFound, nil)
		return
	}
	if r.Method != http.MethodPost {
		customLogger.InfoLogger.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	intent := r.FormValue("intent") // Extract intent (login or signup)
	if intent != "login" && intent != "signup" {
		customLogger.InfoLogger.Println("Invalid intent provided")
		clientError(w, nil, http.StatusBadRequest, nil)
		return
	}

	oauthStateString, err := session.GenerateRandomCSRFText()
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "googleAuthentication", "There is a problem in the process of generating random CSRF text", err))
		serverError(w)
		return
	}
	state := fmt.Sprintf("%s|%s", oauthStateString, intent)
	oauthState[state] = true

	url := googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (handler *HandlerHttp) googleCallback(w http.ResponseWriter, r *http.Request) {

	customLogger.DebugLogger.Println("googleCallback  handler is activated")

	if r.URL.Path != "/oauth2/google/callback" {
		customLogger.InfoLogger.Println("incorrect request's endpoint")
		clientError(w, nil, http.StatusNotFound, nil)
		return
	}

	if r.Method != http.MethodGet {
		customLogger.InfoLogger.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if state == "" || code == "" {
		customLogger.InfoLogger.Println("State or code missing in callback")
		clientError(w, nil, http.StatusBadRequest, nil)
		return
	}

	// Split state to extract CSRF and intent
	parts := strings.Split(state, "|")
	if len(parts) != 2 {
		customLogger.InfoLogger.Println("Invalid state format")
		clientError(w, nil, http.StatusBadRequest, nil)
		return
	}

	intent := parts[1]

	if !oauthState[state] {
		customLogger.InfoLogger.Println("State is invalid, it is corrupted")
		clientError(w, nil, http.StatusBadRequest, nil)
		return
	}
	delete(oauthState, state)

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "googleCallback", "Failed to exchange token", err))
		serverError(w)
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "googleCallback", "Failed to get user info", err))
		serverError(w)
		return
	}
	defer response.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "googleCallback", "Failed to parse user info: "+err.Error(), err))
		serverError(w)
		return
	}

	if intent == "login" {
		tokenSignature, err := handler.Service.LogIn(userInfo.Email, "", "google")
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				customLogger.DebugLogger.Println(fmt.Errorf("Function \"logIn\": %w", err))
				clientError(w, []string{"../ui/html/login.tmpl.html"}, http.StatusBadRequest, domain.ErrUserNotFound)
			} else {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "googleCallback", "There is a problem in the process of giving the tokenSignature", err))
				serverError(w)
			}
			return
		}

		//Cookies
		session.SetTokenToCookie(w, "auth_token", tokenSignature)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if intent == "signup" {
		err = handler.Service.SignUp("", userInfo.Email, "", "google")
		if err != nil {
			if errors.Is(err, domain.ErrInvalidCredential) {
				customLogger.DebugLogger.Println("There is invalid entered Credentials")
				clientError(w, []string{"../ui/html/signup.tmpl.html"}, http.StatusBadRequest, err)
			} else {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "googleCallback", "Failed  Sign up operation", err))
				serverError(w)
			}
			return
		}

		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	}

}
