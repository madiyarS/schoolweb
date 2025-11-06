package models

import "time"

type Document struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	FileType    string    `json:"file_type"`
	Category    string    `json:"category"`
	FolderID    int       `json:"folder_id"`   // Добавлено
	FolderName  string    `json:"folder_name"` // Добавлено
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
