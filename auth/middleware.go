package auth

import (
	"net/http"
	"tiny-invoicing/database"
	"tiny-invoicing/response"
)

// BasicAuth wraps a handler and provides basic authentication.
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			response.Error(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		user, err := database.GetUserByUsername(username)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		if CheckPasswordHash(password, user.PasswordHash) && user.IsAdmin {
			next.ServeHTTP(w, r)
		} else {
			response.Error(w, http.StatusUnauthorized, "Invalid credentials")
		}
	}
}
