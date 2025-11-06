package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"school-website/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func New(filepath string) (*Database, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	database := &Database{db: db}

	if err := database.createTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) createTables() error {
	queries := []string{
		// Существующие таблицы...
		`CREATE TABLE IF NOT EXISTS contacts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT NOT NULL,
            phone TEXT,
            message TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,

		`CREATE TABLE IF NOT EXISTS news (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            image_url TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,

		// Таблица папок
		`CREATE TABLE IF NOT EXISTS folders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            description TEXT,
            icon TEXT DEFAULT 'folder',
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,

		// Обновленная таблица документов
		`CREATE TABLE IF NOT EXISTS documents (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT,
            file_name TEXT NOT NULL,
            file_path TEXT NOT NULL,
            file_size INTEGER NOT NULL,
            file_type TEXT NOT NULL,
            category TEXT,
            folder_id INTEGER,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL
        )`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return err
		}
	}

	log.Println("Database tables created successfully")

	// Выполняем миграции для существующих таблиц
	if err := d.migrateTables(); err != nil {
		return fmt.Errorf("error migrating tables: %v", err)
	}

	// Создаем дефолтные папки
	d.createDefaultFolders()

	return nil
}

// migrateTables выполняет миграции для обновления существующих таблиц
func (d *Database) migrateTables() error {
	// Проверяем и добавляем folder_id в таблицу documents
	if err := d.addColumnIfNotExists("documents", "folder_id", "INTEGER"); err != nil {
		return err
	}

	// Проверяем и добавляем updated_at в таблицу documents
	if err := d.addColumnIfNotExists("documents", "updated_at", "DATETIME DEFAULT CURRENT_TIMESTAMP"); err != nil {
		return err
	}

	return nil
}

// addColumnIfNotExists проверяет существование колонки и добавляет её, если её нет
func (d *Database) addColumnIfNotExists(tableName, columnName, columnDef string) error {
	// Проверяем, существует ли колонка
	var count int
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM pragma_table_info('%s') WHERE name='%s'",
		tableName, columnName,
	)

	err := d.db.QueryRow(query).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking column %s in table %s: %v", columnName, tableName, err)
	}

	// Если колонка не существует, добавляем её
	if count == 0 {
		alterQuery := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, columnName, columnDef)
		if _, err := d.db.Exec(alterQuery); err != nil {
			return fmt.Errorf("error adding column %s to table %s: %v", columnName, tableName, err)
		}
		log.Printf("Added column %s to table %s", columnName, tableName)
	}

	return nil
}

func (d *Database) createDefaultFolders() {
	defaultFolders := []struct {
		name        string
		description string
		icon        string
	}{
		{"Учебные материалы", "Учебные пособия и материалы для уроков", "book"},
		{"Документы школы", "Официальные документы и уставы", "file-text"},
		{"Формы и заявления", "Формы для заполнения родителями", "file-invoice"},
		{"Правила и положения", "Правила внутреннего распорядка", "gavel"},
		{"Расписания", "Расписание занятий и мероприятий", "calendar"},
		{"Прочее", "Другие документы", "folder"},
	}

	for _, folder := range defaultFolders {
		_, err := d.db.Exec(
			`INSERT OR IGNORE INTO folders (name, description, icon) VALUES (?, ?, ?)`,
			folder.name, folder.description, folder.icon,
		)
		if err != nil {
			log.Printf("Warning: failed to create default folder %s: %v", folder.name, err)
		}
	}
}

func (d *Database) Close() error {
	return d.db.Close()
}

// --- Contact Operations ---

func (d *Database) SaveContact(form models.ContactForm) error {
	insertSQL := `INSERT INTO contacts(name, email, phone, message, created_at) VALUES (?, ?, ?, ?, ?)`
	statement, err := d.db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(form.Name, form.Email, form.Phone, form.Message, time.Now())
	if err != nil {
		return fmt.Errorf("error executing insert: %v", err)
	}

	log.Printf("Contact successfully saved: %s (%s)", form.Name, form.Email)
	return nil
}

func (d *Database) GetContacts() ([]models.ContactEntry, error) {
	query := `SELECT id, name, email, COALESCE(phone, '') as phone, message, created_at 
			  FROM contacts ORDER BY created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	var contacts []models.ContactEntry
	for rows.Next() {
		var c models.ContactEntry
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Message, &c.CreatedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		contacts = append(contacts, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %v", err)
	}

	log.Printf("Retrieved %d contacts from database", len(contacts))
	return contacts, nil
}

// --- News Operations ---

func (d *Database) SaveNews(title, content, imageURL string) error {
	insertSQL := `INSERT INTO news(title, content, image_url, created_at) VALUES (?, ?, ?, ?)`
	statement, err := d.db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("error preparing SaveNews statement: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(title, content, imageURL, time.Now())
	if err != nil {
		return fmt.Errorf("error saving news: %v", err)
	}

	log.Printf("News successfully saved: %s", title)
	return nil
}

func (d *Database) GetNews() ([]models.NewsArticle, error) {
	query := `SELECT id, title, content, COALESCE(image_url, '') as image_url, created_at 
			  FROM news ORDER BY created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetNews query failed: %v", err)
	}
	defer rows.Close()

	var articles []models.NewsArticle
	for rows.Next() {
		var a models.NewsArticle
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.ImageURL, &a.CreatedAt); err != nil {
			log.Printf("Error scanning news: %v", err)
			continue
		}
		articles = append(articles, a)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating news: %v", err)
	}

	log.Printf("Retrieved %d news articles from database", len(articles))
	return articles, nil
}

func (d *Database) GetNewsArticle(id string) (models.NewsArticle, error) {
	var a models.NewsArticle
	query := `SELECT id, title, content, COALESCE(image_url, '') as image_url, created_at 
			  FROM news WHERE id = ?`

	err := d.db.QueryRow(query, id).Scan(&a.ID, &a.Title, &a.Content, &a.ImageURL, &a.CreatedAt)
	if err != nil {
		return a, fmt.Errorf("error getting news with ID %s: %v", id, err)
	}

	return a, nil
}

func (d *Database) UpdateNewsArticle(id, title, content, imageURL string) error {
	updateSQL := `UPDATE news SET title = ?, content = ?, image_url = ? WHERE id = ?`
	statement, err := d.db.Prepare(updateSQL)
	if err != nil {
		return fmt.Errorf("error preparing UpdateNewsArticle statement: %v", err)
	}
	defer statement.Close()

	result, err := statement.Exec(title, content, imageURL, id)
	if err != nil {
		return fmt.Errorf("error updating news with ID %s: %v", id, err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Updated %d rows for news with ID %s", rowsAffected, id)

	return nil
}

func (d *Database) DeleteNewsArticle(id string) error {
	log.Printf("Deleting news with ID: %s", id)

	var imageURL string
	err := d.db.QueryRow("SELECT COALESCE(image_url, '') FROM news WHERE id = ?", id).Scan(&imageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("news with ID %s not found", id)
		}
		return fmt.Errorf("error getting news info: %v", err)
	}

	if imageURL != "" && strings.HasPrefix(imageURL, "/uploads/") {
		filePath := "public" + imageURL
		if err := os.Remove(filePath); err != nil {
			log.Printf("Warning: failed to delete file %s: %v", filePath, err)
		} else {
			log.Printf("File %s successfully deleted", filePath)
		}
	}

	deleteSQL := `DELETE FROM news WHERE id = ?`
	result, err := d.db.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("error deleting news: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("news with ID %s not found", id)
	}

	log.Printf("News with ID %s successfully deleted", id)
	return nil
}

// --- Document Operations ---

func (d *Database) SaveDocument(doc models.Document) (int64, error) {
	insertSQL := `INSERT INTO documents(title, description, file_name, file_path, file_size, file_type, category, folder_id, created_at, updated_at) 
                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	statement, err := d.db.Prepare(insertSQL)
	if err != nil {
		return 0, fmt.Errorf("error preparing SaveDocument statement: %v", err)
	}
	defer statement.Close()

	result, err := statement.Exec(
		doc.Title,
		doc.Description,
		doc.FileName,
		doc.FilePath,
		doc.FileSize,
		doc.FileType,
		doc.Category,
		doc.FolderID, // Добавлено
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("error saving document: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %v", err)
	}

	log.Printf("Document successfully saved: %s (ID: %d)", doc.Title, id)
	return id, nil
}

func (d *Database) GetDocuments() ([]models.Document, error) {
	query := `SELECT d.id, d.title, COALESCE(d.description, '') as description, 
              d.file_name, d.file_path, d.file_size, d.file_type, 
              COALESCE(d.category, '') as category, 
              COALESCE(d.folder_id, 0) as folder_id,
              COALESCE(f.name, '') as folder_name,
              d.created_at, d.updated_at
              FROM documents d
              LEFT JOIN folders f ON d.folder_id = f.id
              ORDER BY d.created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetDocuments query failed: %v", err)
	}
	defer rows.Close()

	var documents []models.Document
	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(
			&doc.ID,
			&doc.Title,
			&doc.Description,
			&doc.FileName,
			&doc.FilePath,
			&doc.FileSize,
			&doc.FileType,
			&doc.Category,
			&doc.FolderID,
			&doc.FolderName,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning document: %v", err)
			continue
		}
		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %v", err)
	}

	log.Printf("Retrieved %d documents from database", len(documents))
	return documents, nil
}

func (d *Database) GetDocument(id string) (models.Document, error) {
	var doc models.Document
	query := `SELECT d.id, d.title, COALESCE(d.description, '') as description, d.file_name, d.file_path,
			  d.file_size, d.file_type, COALESCE(d.category, '') as category, 
			  COALESCE(d.folder_id, 0) as folder_id,
			  COALESCE(f.name, '') as folder_name,
			  d.created_at, d.updated_at
			  FROM documents d
			  LEFT JOIN folders f ON d.folder_id = f.id
			  WHERE d.id = ?`

	err := d.db.QueryRow(query, id).Scan(
		&doc.ID,
		&doc.Title,
		&doc.Description,
		&doc.FileName,
		&doc.FilePath,
		&doc.FileSize,
		&doc.FileType,
		&doc.Category,
		&doc.FolderID,
		&doc.FolderName,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return doc, fmt.Errorf("document with ID %s not found", id)
		}
		return doc, fmt.Errorf("error getting document with ID %s: %v", id, err)
	}

	return doc, nil
}

func (d *Database) GetDocumentsByCategory(category string) ([]models.Document, error) {
	query := `SELECT d.id, d.title, COALESCE(d.description, '') as description, d.file_name, d.file_path,
			  d.file_size, d.file_type, COALESCE(d.category, '') as category, 
			  COALESCE(d.folder_id, 0) as folder_id,
			  COALESCE(f.name, '') as folder_name,
			  d.created_at, d.updated_at
			  FROM documents d
			  LEFT JOIN folders f ON d.folder_id = f.id
			  WHERE d.category = ? ORDER BY d.created_at DESC`

	rows, err := d.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("GetDocumentsByCategory query failed: %v", err)
	}
	defer rows.Close()

	var documents []models.Document
	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(
			&doc.ID,
			&doc.Title,
			&doc.Description,
			&doc.FileName,
			&doc.FilePath,
			&doc.FileSize,
			&doc.FileType,
			&doc.Category,
			&doc.FolderID,
			&doc.FolderName,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning document: %v", err)
			continue
		}
		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %v", err)
	}

	log.Printf("Retrieved %d documents for category '%s'", len(documents), category)
	return documents, nil
}

func (d *Database) DeleteDocument(id string) error {
	log.Printf("Deleting document with ID: %s", id)

	// Get document info first to delete the file
	var filePath string
	err := d.db.QueryRow("SELECT file_path FROM documents WHERE id = ?", id).Scan(&filePath)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("document with ID %s not found", id)
		}
		return fmt.Errorf("error getting document info: %v", err)
	}

	// Delete the file from filesystem
	if filePath != "" {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Warning: failed to delete file %s: %v", filePath, err)
		} else {
			log.Printf("File %s successfully deleted", filePath)
		}
	}

	// Delete from database
	deleteSQL := `DELETE FROM documents WHERE id = ?`
	result, err := d.db.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("error deleting document: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document with ID %s not found", id)
	}

	log.Printf("Document with ID %s successfully deleted", id)
	return nil
}

// --- Folder Operations ---

func (d *Database) GetFolders() ([]models.Folder, error) {
	query := `SELECT id, name, COALESCE(description, '') as description, 
			  COALESCE(icon, 'folder') as icon, created_at 
			  FROM folders ORDER BY name ASC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetFolders query failed: %v", err)
	}
	defer rows.Close()

	var folders []models.Folder
	for rows.Next() {
		var folder models.Folder
		if err := rows.Scan(&folder.ID, &folder.Name, &folder.Description, &folder.Icon, &folder.CreatedAt); err != nil {
			log.Printf("Error scanning folder: %v", err)
			continue
		}
		folders = append(folders, folder)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating folders: %v", err)
	}

	log.Printf("Retrieved %d folders from database", len(folders))
	return folders, nil
}

func (d *Database) CreateFolder(name, description, icon string) (int64, error) {
	insertSQL := `INSERT INTO folders(name, description, icon, created_at) VALUES (?, ?, ?, ?)`
	statement, err := d.db.Prepare(insertSQL)
	if err != nil {
		return 0, fmt.Errorf("error preparing CreateFolder statement: %v", err)
	}
	defer statement.Close()

	result, err := statement.Exec(name, description, icon, time.Now())
	if err != nil {
		return 0, fmt.Errorf("error creating folder: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %v", err)
	}

	log.Printf("Folder successfully created: %s (ID: %d)", name, id)
	return id, nil
}

func (d *Database) DeleteFolder(id string) error {
	log.Printf("Deleting folder with ID: %s", id)

	// Check if folder exists
	var folderName string
	err := d.db.QueryRow("SELECT name FROM folders WHERE id = ?", id).Scan(&folderName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("folder with ID %s not found", id)
		}
		return fmt.Errorf("error getting folder info: %v", err)
	}

	// Delete the folder (documents will have folder_id set to NULL due to ON DELETE SET NULL)
	deleteSQL := `DELETE FROM folders WHERE id = ?`
	result, err := d.db.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("error deleting folder: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("folder with ID %s not found", id)
	}

	log.Printf("Folder with ID %s (%s) successfully deleted", id, folderName)
	return nil
}

func (d *Database) GetDocumentsByFolder(folderID string) ([]models.Document, error) {
	query := `SELECT d.id, d.title, COALESCE(d.description, '') as description, 
              d.file_name, d.file_path, d.file_size, d.file_type, 
              COALESCE(d.category, '') as category, 
              COALESCE(d.folder_id, 0) as folder_id,
              COALESCE(f.name, '') as folder_name,
              d.created_at, d.updated_at
              FROM documents d
              LEFT JOIN folders f ON d.folder_id = f.id
              WHERE d.folder_id = ?
              ORDER BY d.created_at DESC`

	rows, err := d.db.Query(query, folderID)
	if err != nil {
		return nil, fmt.Errorf("GetDocumentsByFolder query failed: %v", err)
	}
	defer rows.Close()

	var documents []models.Document
	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(
			&doc.ID,
			&doc.Title,
			&doc.Description,
			&doc.FileName,
			&doc.FilePath,
			&doc.FileSize,
			&doc.FileType,
			&doc.Category,
			&doc.FolderID,
			&doc.FolderName,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning document: %v", err)
			continue
		}
		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %v", err)
	}

	log.Printf("Retrieved %d documents for folder ID '%s'", len(documents), folderID)
	return documents, nil
}

func (d *Database) Ping() error {
	return d.db.Ping()
}
