package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"school-website/internal/database"
	"school-website/internal/models"
)

type ContactHandler struct {
	db *database.Database
}

func NewContactHandler(db *database.Database) *ContactHandler {
	return &ContactHandler{db: db}
}

func (h *ContactHandler) SubmitContact(w http.ResponseWriter, r *http.Request) {
	// Allow CORS for frontend requests (for development)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var form models.ContactForm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&form)
	if err != nil {
		log.Printf("Error decoding contact form: %v", err)
		http.Error(w, "Error reading data", http.StatusBadRequest)
		return
	}

	// Save data to database
	err = h.db.SaveContact(form)
	if err != nil {
		log.Printf("Error saving contact: %v", err)
		http.Error(w, "Error saving data to database", http.StatusInternalServerError)
		return
	}

	log.Printf("Message from %s (%s) successfully saved to database", form.Name, form.Email)

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Message successfully received and saved!",
	})
}

func (h *ContactHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== Processing request /admin/api/applications ===")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contacts, err := h.db.GetContacts()
	if err != nil {
		log.Printf("Error getting contacts: %v", err)
		http.Error(w, "Failed to get contacts", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully retrieved %d contacts", len(contacts))
	json.NewEncoder(w).Encode(contacts)
}
