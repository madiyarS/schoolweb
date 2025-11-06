package handlers

import (
	"encoding/json"
	"net/http"

	"school-website/internal/database"
	"school-website/internal/models"

	"github.com/gorilla/mux"
)

type FolderHandler struct {
	db *database.Database
}

func NewFolderHandler(db *database.Database) *FolderHandler {
	return &FolderHandler{db: db}
}

func (h *FolderHandler) GetAllFolders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	folders, err := h.db.GetFolders()
	if err != nil {
		http.Error(w, "Failed to get folders", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(folders)
}

func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var folder models.Folder
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if folder.Name == "" {
		http.Error(w, "Folder name is required", http.StatusBadRequest)
		return
	}

	if folder.Icon == "" {
		folder.Icon = "folder"
	}

	id, err := h.db.CreateFolder(folder.Name, folder.Description, folder.Icon)
	if err != nil {
		http.Error(w, "Failed to create folder", http.StatusInternalServerError)
		return
	}

	folder.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(folder)
}

func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.db.DeleteFolder(id); err != nil {
		http.Error(w, "Failed to delete folder", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Folder deleted successfully"})
}

func (h *FolderHandler) GetFolderDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	folderID := vars["id"]

	documents, err := h.db.GetDocumentsByFolder(folderID)
	if err != nil {
		http.Error(w, "Failed to get documents", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(documents)
}
