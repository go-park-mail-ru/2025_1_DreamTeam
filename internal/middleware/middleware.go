package middleware

import (
	"net/http"
)

func CorsOptionsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://217.16.21.64")
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8001")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Credentials")

		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
		}

	})
}
