// internal/infrastructure/middleware/auth.go
package middleware

import "net/http"

// AuthMiddleware is a stub for future authentication middleware.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
