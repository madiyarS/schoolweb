package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

type AuthMiddleware struct {
	store *sessions.CookieStore
}

func NewAuthMiddleware(store *sessions.CookieStore) *AuthMiddleware {
	return &AuthMiddleware{store: store}
}

func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := am.store.Get(r, "session-name")

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// For API requests, return 401 instead of redirect
			if strings.Contains(r.URL.Path, "/api/") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Redirect to login page if not authenticated
			http.Redirect(w, r, "/admin/login.html", http.StatusFound)
			return
		}

		// If authenticated, pass control to the next handler
		next.ServeHTTP(w, r)
	})
}
