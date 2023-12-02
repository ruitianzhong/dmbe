package authentication

import (
	"log"
	"net/http"
)

// AuthMiddleware just do nothing for now
func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/auth/login" {
				session, _ := store.Get(r, "dm-session")
				if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
					log.Println("reject")
				} else {
					log.Println("authenticated", auth)
				}
				next.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})

}
