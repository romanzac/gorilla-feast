package middleware

import (
	"github.com/gorilla/csrf"
	"net/http"
)

func CSRFTokenToHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		next.ServeHTTP(w, r)
	})
}
