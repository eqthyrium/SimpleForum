package middleware

import (
	"SimpleForum/internal/service/auth"
	"SimpleForum/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"time"
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
*/
var customLogger *logger.Logger = logger.NewLogger().GetLoggerObject("../../../../logging/info.log", "../../../../logging/error.log", "Middleware")

func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)
		customLogger.InfoLogger.Print(fmt.Sprintf("[%s] %s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr, duration))
	})
}

func RoleAdjusterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetTokenFromCookie(r)
		if err != nil {
			// InternalServer Error
			return
		}
		if tokenString == "" {
			ctx := context.WithValue(r.Context(), "Role", "Guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		verification, err := auth.VerifyToken(tokenString)

		if err != nil {
			// InternalServer Error
			return
		}

		if !verification {
			ctx := context.WithValue(r.Context(), "Role", "Guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		extractedToken, err := auth.ExtractDataFromToken(tokenString)
		if err != nil {
			// InternalServer Error
			return
		}
		if auth.MapUUID[extractedToken.UserId] != extractedToken.UUID {
			ctx := context.WithValue(r.Context(), "Role", "Guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		switch auth.CheckTokenTime(extractedToken) {
		case "Invalid-Token":
			ctx := context.WithValue(r.Context(), "Role", "Guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		case "Extend-Token":
			extendedToken, err := auth.ExtendTokenExistence(extractedToken)
			if err != nil {
				// InternalServer Error
				ctx := context.WithValue(r.Context(), "Role", "Guest")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			auth.SetTokenToCookie(w, extendedToken)
		}

		ctx := context.WithValue(r.Context(), "Role", extractedToken.Role)
		ctx = context.WithValue(r.Context(), "UserId", extractedToken.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
