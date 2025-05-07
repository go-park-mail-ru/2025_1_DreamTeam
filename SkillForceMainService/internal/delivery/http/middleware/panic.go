package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

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
