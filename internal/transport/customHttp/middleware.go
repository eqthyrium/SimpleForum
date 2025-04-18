package customHttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"SimpleForum/internal/domain"
	session2 "SimpleForum/internal/transport/session"
	"SimpleForum/pkg/logger"
)

/*
ToDo

1)Security Middleware (done)

2) Logging Middleware (done)

3) Token Validation Middleware (done)
Check presence of token in cookie from client request
├── True
│   ├── Check for verification of token
│   │   ├── True
│   │   │   ├── Check UserId existence in the MapUUID
│   │   │   │   ├── True
│   │   │   │   │   ├── Check expiration of the token
│   │   │   │   │   │   ├── True
│   │   │   │   │   │   │   ├── Check if token's time surpasses threshold time
│   │   │   │   │   │   │   │   ├── True
│   │   │   │   │   │   │   │   │   ├── Extend token time by 45 minutes and send new token in cookie to client
│   │   │   │   │   │   │   │   └── False
│   │   │   │   │   │   │   │       ├── Send appropriate webpage
│   │   │   │   │   │   └── False
│   │   │   │   │   │       ├── Send guest homepage (token expired)
│   │   │   │   └── False
│   │   │   │       ├── Send guest homepage (UserId not in MapUUID or another UUID found)
│   │   └── False
│   │       ├── Send guest webpage (failed token verification)
└── False
├── Send guest webpage

4) I have to implement CSRF checking middleware

5) Panic middleware
*/

var CSRFMap map[string]string = make(map[string]string)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		customLogger.InfoLogger.Print(fmt.Sprintf("Method:[%v], URL_Path: %v, Remote_Address: %v\n", r.Method, r.URL.Path, r.RemoteAddr))
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)
		// Here must be reponse's status code
		customLogger.InfoLogger.Print(fmt.Sprintf("The End of the client request, and its Duration:%v\n", duration))
	})
}

func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

func RoleAdjusterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customLogger.DebugLogger.Println("The RoleAdjusterMiddleware is started")

		tokenString, err := session2.GetTokenFromCookie(r, "auth_token")
		if err != nil {
			customLogger.DebugLogger.Println("There is an error about getting the token from the cookie!!!")
			if errors.Is(err, http.ErrNoCookie) {
				customLogger.InfoLogger.Println("There is no cookie in the request of the client")
				ctx := context.WithValue(r.Context(), "Role", "Guest")
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				customLogger.InfoLogger.Println("There is a problem in the process of Extraction token from cookie")
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "RoleAdjusterMiddleware", "There is a problem in the process of Extraction token from cookie", err))
				session2.DeleteSessionCookie(w, "auth_token")
				serverError(w)
			}
			return
		}

		if tokenString == "" {
			customLogger.DebugLogger.Println("Entered into the absence of the token in the cookie")
			customLogger.InfoLogger.Println("There is cookie, but there is no token")
			session2.DeleteSessionCookie(w, "auth_token")
			ctx := context.WithValue(r.Context(), "Role", "Guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		err = session2.VerifyToken(tokenString)
		if err != nil {
			customLogger.DebugLogger.Println("Entered into error handling of verification of the token")
			session2.DeleteSessionCookie(w, "auth_token")

			if errors.Is(err, domain.ErrInvalidToken) {
				customLogger.InfoLogger.Println("There is an invalid token")
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			} else {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "RoleAdjusterMiddleware", "There is a problem in the process of verification of the token", err))
				serverError(w)
			}
			return
		}

		extractedToken, err := session2.ExtractDataFromToken(tokenString)
		if err != nil {
			customLogger.DebugLogger.Println("Entered into error handling of extraction token")
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "RoleAdjusterMiddleware", "There is a problem in the process of extraction of the token", err))
			session2.DeleteSessionCookie(w, "auth_token")
			serverError(w)
			return
		}

		if session2.MapUUID[extractedToken.UserId] != extractedToken.UUID {
			customLogger.DebugLogger.Println("Entered into error handling of the check up of MappUUID")
			customLogger.InfoLogger.Println("There is not current token for the client")
			session2.DeleteSessionCookie(w, "auth_token")
			// delete(CSRFMap, session2.MapUUID[extractedToken.UserId])
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		customLogger.DebugLogger.Println("Checking part of token's time")
		switch session2.CheckTokenTime(extractedToken) {
		case "Expired-Token":
			customLogger.InfoLogger.Println("There is an expired token")
			session2.DeleteSessionCookie(w, "auth_token")
			delete(CSRFMap, session2.MapUUID[extractedToken.UserId])
			delete(session2.MapUUID, extractedToken.UserId)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		case "Extend-Token":
			extendedToken, err := session2.ExtendTokenExistence(extractedToken)
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "RoleAdjusterMiddleware", "There is a problem in the process of extension of time of the token", err))
				session2.DeleteSessionCookie(w, "auth_token")
				serverError(w)
				return
			}
			customLogger.InfoLogger.Println("The member with userId:", extractedToken.UserId, "and its role:", extractedToken.Role, ", its expireTime is refreshed by adding 45 min to previous left time.")
			session2.SetTokenToCookie(w, "auth_token", extendedToken)
		}

		ctx := context.WithValue(r.Context(), "Role", extractedToken.Role)
		ctx = context.WithValue(ctx, "UserId", extractedToken.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customLogger.DebugLogger.Println("The CSRFMiddleware is started")

		role := r.Context().Value("Role").(string)

		if role != "Guest" && r.Method == http.MethodPost {
			customLogger.DebugLogger.Println("inside of the CSRFMiddleware's Post method checking part ")

			formCSRFText := r.FormValue("csrf_text")
			userId := r.Context().Value("UserId").(int)

			if formCSRFText != CSRFMap[session2.MapUUID[userId]] {
				// fmt.Println("Fucking CSRF attack")
				// fmt.Println("MapUUID[userId]:", session2.MapUUID[userId])
				// fmt.Println("HTML formCSRFTEXT", formCSRFText, "\n", "Server Map:", CSRFMap[session2.MapUUID[userId]])
				customLogger.InfoLogger.Println("The CSRF attack is detected, its IP is:", r.RemoteAddr)
				//if _, ok := CSRFMap[session2.MapUUID[userId]]; ok {
				//	delete(CSRFMap, session2.MapUUID[userId])
				//}
				//delete(session2.MapUUID, userId)
				session2.DeleteSessionCookie(w, "auth_token")
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}
			customLogger.DebugLogger.Println("The CSRF checking part was good")

		}

		next.ServeHTTP(w, r)
	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				customLogger.ErrorLogger.Println("Panic:\n", err, string(debug.Stack()))
				serverError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RateLimiterMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	tokens := limit
	resetTime := time.Now().Add(window)
	var mu sync.Mutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			defer mu.Unlock()

			customLogger.InfoLogger.Println("Tokens left: %d, Reset time: %v\n", tokens, resetTime)

			if time.Now().After(resetTime) {
				tokens = limit
				resetTime = time.Now().Add(window)
				customLogger.InfoLogger.Println("Resetting tokens.")

			}

			if tokens <= 0 {
				customLogger.ErrorLogger.Println("Rate limit exceeded.")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			tokens--
			next.ServeHTTP(w, r)
		})
	}
}
