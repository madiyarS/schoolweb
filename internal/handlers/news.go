package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"school-website/internal/database"
	"school-website/internal/services"

	"github.com/gorilla/mux"
)

type NewsHandler struct {
	db            *database.Database
	uploadService *services.FileUploadService
}

func NewNewsHandler(db *database.Database, uploadService *services.FileUploadService) *NewsHandler {
	return &NewsHandler{
		db:            db,
		uploadService: uploadService,
	}
}

func (h *NewsHandler) GetAllNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	articles, err := h.db.GetNews()
	if err != nil {
		log.Printf("Error getting news: %v", err)
		http.Error(w, "Failed to get news", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func (h *NewsHandler) GetSingleNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing news ID", http.StatusBadRequest)
		return
	}

	article, err := h.db.GetNewsArticle(id)
	if err != nil {
		log.Printf("Error getting news with ID %s: %v", id, err)
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func (h *NewsHandler) CreateNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	imageURLFromForm := strings.TrimSpace(r.FormValue("image_url"))

	log.Printf("Creating news: title=%s, content length=%d, imageURLFromForm=%s",
		title, len(content), imageURLFromForm)

	if title == "" || content == "" {
		log.Printf("Missing required fields: title='%s', content length=%d", title, len(content))
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Priority for uploaded file
	finalImageURL, err := h.uploadService.HandleFileUpload(r)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
		return
	}

	if finalImageURL == "" {
		finalImageURL = imageURLFromForm
		log.Printf("Using URL from form: %s", finalImageURL)
	}

	err = h.db.SaveNews(title, content, finalImageURL)
	if err != nil {
		log.Printf("Error saving news to database: %v", err)
		http.Error(w, "Failed to save news", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":   true,
		"message":   "News successfully created",
		"image_url": finalImageURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *NewsHandler) UpdateNews(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Error parsing form during update: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing news ID", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	imageURLFromForm := strings.TrimSpace(r.FormValue("image_url"))

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	existingArticle, err := h.db.GetNewsArticle(id)
	if err != nil {
		log.Printf("Article with ID %s not found: %v", id, err)
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	newImageURL, err := h.uploadService.HandleFileUpload(r)
	if err != nil {
		log.Printf("Error uploading file during update: %v", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	finalImageURL := existingArticle.ImageURL

	if newImageURL != "" {
		finalImageURL = newImageURL
	} else if imageURLFromForm != existingArticle.ImageURL {
		finalImageURL = imageURLFromForm
	}

	err = h.db.UpdateNewsArticle(id, title, content, finalImageURL)
	if err != nil {
		log.Printf("Error updating news: %v", err)
		http.Error(w, "Failed to update news", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "News updated successfully"})
}

func (h *NewsHandler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Printf("Missing news ID in delete request")
		http.Error(w, "Missing news ID", http.StatusBadRequest)
		return
	}

	log.Printf("Delete request for news ID: %s", id)

	err := h.db.DeleteNewsArticle(id)
	if err != nil {
		log.Printf("Error deleting news with ID %s: %v", id, err)
		if strings.Contains(err.Error(), "не найдена") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete news article", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("News with ID %s successfully deleted", id)

	response := map[string]interface{}{
		"success": true,
		"message": "News successfully deleted",
		"id":      id,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
