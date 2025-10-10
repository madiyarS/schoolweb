package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileUploadService struct {
	uploadDir string
}

func NewFileUploadService(uploadDir string) *FileUploadService {
	return &FileUploadService{uploadDir: uploadDir}
}

func (s *FileUploadService) HandleFileUpload(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil // File was not uploaded, not an error
		}
		log.Printf("Error getting file from form: %v", err)
		return "", err
	}
	defer file.Close()

	// Validate file extension (basic validation)
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	// Create unique filename
	uniqueFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	uniqueFileName = strings.ReplaceAll(uniqueFileName, " ", "_")

	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		log.Printf("Error creating upload directory: %v", err)
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	filePath := filepath.Join(s.uploadDir, uniqueFileName)
	log.Printf("Attempting to save file: %s", filePath)

	// Create file on server
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file %s: %v", filePath, err)
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy uploaded file content
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Error copying file content: %v", err)
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	finalURL := "/uploads/" + uniqueFileName
	log.Printf("File successfully saved: %s -> URL: %s", filePath, finalURL)
	return finalURL, nil
}
