package customHttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"SimpleForum/internal/config"
	"SimpleForum/internal/domain"
	"SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
)

func (handler *HandlerHttp) githubAuthentication(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("githubAuthentication handler is activated")

	if r.URL.Path != "/oauth2/github" {
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
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "githubAuthentication", "There is a problem in the process of generating random CSRF text", err))
		serverError(w)
		return
	}
	state := fmt.Sprintf("%s|%s", oauthStateString, intent)
	oauthState[state] = true
	http.Redirect(w, r, fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s\n",
		config.Config.GithubOauth.ClientID, config.Config.GithubOauth.RedirectURI, state), http.StatusFound)
}

func (handler *HandlerHttp) githubCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/oauth2/github/callback" {
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

	// Exchange code for access token
	tokenURL := "https://github.com/login/oauth/access_token"
	data := url.Values{
		"client_id":     {config.Config.GithubOauth.ClientID},
		"client_secret": {config.Config.GithubOauth.ClientSecret},
		"code":          {code},
		"redirect_uri":  {config.Config.GithubOauth.RedirectURI},
	}
	req, _ := http.NewRequest("POST", tokenURL, nil)
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		customLogger.ErrorLogger.Print("Failed to get token")
		serverError(w)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var tokenResponse map[string]interface{}
	json.Unmarshal(body, &tokenResponse)
	accessToken := tokenResponse["access_token"].(string)
	emailReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	emailReq.Header.Set("Authorization", "Bearer "+accessToken)
	emailReq.Header.Set("Accept", "application/vnd.github.v3+json")
	emailResp, err := client.Do(emailReq)
	if err != nil {
		customLogger.ErrorLogger.Print("Failed to fetch email")
		serverError(w)
		return
	}
	defer emailResp.Body.Close()

	if emailResp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(emailResp.Body)
		customLogger.ErrorLogger.Print("Failed to fetch email" + fmt.Sprintf("Failed to fetch email. Status: %d, Body: %s\n", emailResp.StatusCode, string(body)))
		serverError(w)
		return
	}

	emailBody, _ := ioutil.ReadAll(emailResp.Body)
	customLogger.DebugLogger.Print(fmt.Sprintf("Emails Response Body:%s", string(emailBody)))

	var emails []map[string]interface{}
	if err := json.Unmarshal(emailBody, &emails); err != nil {
		customLogger.ErrorLogger.Println("Failed to parse email response")
		serverError(w)
		return
	}

	if len(emails) == 0 {
		customLogger.ErrorLogger.Println("No emails are found")
		serverError(w)
		return
	}

	// Extract the primary email
	primaryEmail := ""
	for _, email := range emails {
		if email["primary"].(bool) {
			primaryEmail = email["email"].(string)
			break
		}
	}

	if primaryEmail == "" {
		customLogger.ErrorLogger.Println("Primary email not found")
		serverError(w)
		return
	}

	if intent == "login" {
		tokenSignature, err := handler.Service.LogIn(primaryEmail, "", "google")
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

		// Cookies
		session.SetTokenToCookie(w, "auth_token", tokenSignature)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if intent == "signup" {
		err = handler.Service.SignUp("", primaryEmail, "", "google")
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
