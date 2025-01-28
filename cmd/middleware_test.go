package main

import (
	"SimpleForum/internal/transport/customHttp"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterMiddleware(t *testing.T) {
	// Лимитируем до 2 запросов за 3 секунды
	limit := 2
	window := 3 * time.Second
	rateLimiter := customHttp.RateLimiterMiddleware(limit, window)

	handler := rateLimiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 4; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if i < limit && rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d on request %d", rec.Code, i+1)
		} else if i >= limit && rec.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status 429, got %d on request %d", rec.Code, i+1)
		}
	}
}
