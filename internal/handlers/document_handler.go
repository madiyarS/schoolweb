package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
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

	// Parse multipart form (500 MB max for multiple large files)
	// This allows uploading multiple files at once
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":    "Request too large",
			"message":  "The total size of uploaded files exceeds the maximum allowed size (500MB). Please try uploading fewer files or smaller files.",
			"max_size": "500MB",
		})
		return
	}

	// Get form values
	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")
	folderIDStr := r.FormValue("folder_id")

	// Get all files from form
	files := r.MultipartForm.File["document"]

	if len(files) == 0 {
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	// Parse folder ID
	var folderID int
	if folderIDStr != "" {
		folderID, _ = strconv.Atoi(folderIDStr)
	}

	// Handle single file (backward compatibility)
	if len(files) == 1 {
		// Use filename (without extension) as title if title is empty
		fileTitle := title
		if fileTitle == "" {
			filename := files[0].Filename
			ext := filepath.Ext(filename)
			if ext != "" {
				fileTitle = filename[:len(filename)-len(ext)]
			} else {
				fileTitle = filename
			}
		}

		file, err := files[0].Open()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			http.Error(w, "Failed to open file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		doc, err := h.service.UploadDocument(fileTitle, description, category, folderID, file, files[0])
		if err != nil {
			log.Printf("Error uploading document: %v", err)
			http.Error(w, fmt.Sprintf("Failed to upload document: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(doc)
		log.Printf("Document uploaded successfully: %s (ID: %d)", doc.Title, doc.ID)
		return
	}

	// Handle multiple files - use filename as title for each if title is empty
	documents, errors := h.service.UploadMultipleDocuments(title, description, category, folderID, files)

	response := map[string]interface{}{
		"success":   len(documents),
		"failed":    len(errors),
		"documents": documents,
	}

	if len(errors) > 0 {
		errorMessages := make([]string, len(errors))
		for i, err := range errors {
			errorMessages[i] = err.Error()
		}
		response["errors"] = errorMessages
		log.Printf("Uploaded %d documents, %d failed", len(documents), len(errors))
	} else {
		log.Printf("Successfully uploaded %d documents", len(documents))
	}

	w.Header().Set("Content-Type", "application/json")
	if len(errors) > 0 && len(documents) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	} else if len(errors) > 0 {
		w.WriteHeader(http.StatusPartialContent) // 206 for partial success
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	json.NewEncoder(w).Encode(response)
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
