package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"school-website/internal/services"

	"github.com/gorilla/mux"
)

type DocumentHandler struct {
	service *services.DocumentService
}

func NewDocumentHandler(service *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse multipart form (32 MB max)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("document")
	if err != nil {
		log.Printf("Error getting file: %v", err)
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get form values
	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")
	folderIDStr := r.FormValue("folder_id") // Добавлено

	// Validate required fields
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Parse folder ID
	var folderID int
	if folderIDStr != "" {
		folderID, _ = strconv.Atoi(folderIDStr)
	}

	// Upload document
	doc, err := h.service.UploadDocument(title, description, category, folderID, file, fileHeader)
	if err != nil {
		log.Printf("Error uploading document: %v", err)
		http.Error(w, fmt.Sprintf("Failed to upload document: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
	log.Printf("Document uploaded successfully: %s (ID: %d)", doc.Title, doc.ID)
}

func (h *DocumentHandler) GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	category := r.URL.Query().Get("category")

	var documents interface{}
	var err error

	if category != "" {
		documents, err = h.service.GetDocumentsByCategory(category)
	} else {
		documents, err = h.service.GetAllDocuments()
	}

	if err != nil {
		log.Printf("Error getting documents: %v", err)
		http.Error(w, "Failed to get documents", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(documents)
}

func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	doc, err := h.service.GetDocument(id)
	if err != nil {
		log.Printf("Error getting document: %v", err)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(doc)
}

func (h *DocumentHandler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	doc, err := h.service.GetDocument(id)
	if err != nil {
		log.Printf("Error getting document for download: %v", err)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Set headers for download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", doc.FileName))
	w.Header().Set("Content-Type", doc.FileType)
	w.Header().Set("Content-Length", strconv.FormatInt(doc.FileSize, 10))

	// Serve the file
	http.ServeFile(w, r, doc.FilePath)
	log.Printf("Document downloaded: %s (ID: %d)", doc.FileName, doc.ID)
}

func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteDocument(id); err != nil {
		log.Printf("Error deleting document: %v", err)
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Document deleted successfully"}`))
	log.Printf("Document deleted successfully: ID %s", id)
}
