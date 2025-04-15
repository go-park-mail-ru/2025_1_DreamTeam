package middleware

import (
	"net/http"
	"skillForce/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	csrfHeader = "X-CSRF-Token"
	csrfExpiry = 40 * time.Minute
)

var jwtSecret []byte

func InitCSRF(cfg *config.Config) {
	jwtSecret = []byte(cfg.Secrets.JwtSessionSecret)
}

func GenerateCSRFToken() (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(csrfExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			csrfToken, err := GenerateCSRFToken()
			if err != nil {
				http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
				return
			}
			w.Header().Set(csrfHeader, csrfToken)
			next.ServeHTTP(w, r)
			return
		}

		clientToken := r.Header.Get(csrfHeader)
		if clientToken == "" {
			http.Error(w, "Missing CSRF token", http.StatusForbidden)
			return
		}

		csrfToken, err := jwt.Parse(clientToken, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !csrfToken.Valid {
			http.Error(w, "Invalid or expired CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
