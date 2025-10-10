package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"school-website/internal/config"
	"school-website/internal/models"

	"github.com/gorilla/sessions"
)

type AuthHandler struct {
	store  *sessions.CookieStore
	config *config.Config
}

func NewAuthHandler(store *sessions.CookieStore, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		store:  store,
		config: cfg,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Printf("Error decoding login data: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt for user: %s", creds.Username)

	// Check credentials
	if creds.Username != h.config.AdminUsername || creds.Password != h.config.AdminPassword {
		log.Printf("Invalid credentials for user: %s", creds.Username)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create session
	session, _ := h.store.Get(r, "session-name")
	session.Values["authenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("User %s successfully authenticated", creds.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")

	// Reset authentication flag
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1 // Delete cookie

	err := session.Save(r, w)
	if err != nil {
		log.Printf("Error during logout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("User logged out")
	http.Redirect(w, r, "/admin/login.html", http.StatusFound)
}
