package authentication

import (
	"net/http"
)

// AuthMiddleware just do nothing for now
func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Session-Token")
			if token != "42" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		})

}
