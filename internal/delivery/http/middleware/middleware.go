package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

func CorsOptionsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://217.16.21.64")
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8001")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Credentials, X-CSRF-Token")
		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")

		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
		}

	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				log.Print("panicMiddleware ", r.URL.Path)
				log.Print("recovered ", err)
				debug.PrintStack()
				http.Error(w, "server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
