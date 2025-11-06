package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"school-website/internal/database"
	"school-website/internal/models"
)

type DocumentService struct {
	db         *database.Database
	uploadPath string
}

func NewDocumentService(db *database.Database, uploadPath string) *DocumentService {
	// Ensure upload directory exists
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		fmt.Printf("Warning: failed to create upload directory %s: %v\n", uploadPath, err)
	}

	return &DocumentService{
		db:         db,
		uploadPath: uploadPath,
	}
}

func (s *DocumentService) UploadDocument(title, description, category string, folderID int, file multipart.File, fileHeader *multipart.FileHeader) (*models.Document, error) {
	// Generate unique filename
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s", timestamp, fileHeader.Filename)
	filePath := filepath.Join(s.uploadPath, fileName)

	// Create file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Create document model
	doc := models.Document{
		Title:       title,
		Description: description,
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    fileHeader.Size,
		FileType:    fileHeader.Header.Get("Content-Type"),
		Category:    category,
		FolderID:    folderID, // Добавлено
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	id, err := s.db.SaveDocument(doc)
	if err != nil {
		os.Remove(filePath) // Clean up file if database save fails
		return nil, fmt.Errorf("failed to save document to database: %v", err)
	}

	doc.ID = int(id)
	return &doc, nil
}

func (s *DocumentService) GetAllDocuments() ([]models.Document, error) {
	return s.db.GetDocuments()
}

func (s *DocumentService) GetDocument(id string) (*models.Document, error) {
	doc, err := s.db.GetDocument(id)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *DocumentService) GetDocumentsByCategory(category string) ([]models.Document, error) {
	return s.db.GetDocumentsByCategory(category)
}

func (s *DocumentService) DeleteDocument(id string) error {
	return s.db.DeleteDocument(id)
}
