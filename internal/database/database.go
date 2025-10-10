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
	createContactsTableSQL := `CREATE TABLE IF NOT EXISTS contacts (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT NOT NULL,
		"email" TEXT NOT NULL,
        "phone" TEXT,
		"message" TEXT NOT NULL,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := d.db.Exec(createContactsTableSQL); err != nil {
		return fmt.Errorf("error creating contacts table: %v", err)
	}
	log.Println("Table 'contacts' successfully created or already exists")

	createNewsTableSQL := `CREATE TABLE IF NOT EXISTS news (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT NOT NULL,
		"content" TEXT NOT NULL,
        "image_url" TEXT,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := d.db.Exec(createNewsTableSQL); err != nil {
		return fmt.Errorf("error creating news table: %v", err)
	}
	log.Println("Table 'news' successfully created or already exists")

	// Create documents table
	createDocumentsTableSQL := `CREATE TABLE IF NOT EXISTS documents (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT NOT NULL,
		"description" TEXT,
		"file_name" TEXT NOT NULL,
		"file_path" TEXT NOT NULL,
		"file_size" INTEGER NOT NULL,
		"file_type" TEXT NOT NULL,
		"category" TEXT,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP,
		"updated_at" DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := d.db.Exec(createDocumentsTableSQL); err != nil {
		return fmt.Errorf("error creating documents table: %v", err)
	}
	log.Println("Table 'documents' successfully created or already exists")

	return nil
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
	insertSQL := `INSERT INTO documents(title, description, file_name, file_path, file_size, file_type, category, created_at, updated_at) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

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
	query := `SELECT id, title, COALESCE(description, '') as description, file_name, file_path, 
			  file_size, file_type, COALESCE(category, '') as category, created_at, updated_at
			  FROM documents ORDER BY created_at DESC`

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
	query := `SELECT id, title, COALESCE(description, '') as description, file_name, file_path,
			  file_size, file_type, COALESCE(category, '') as category, created_at, updated_at
			  FROM documents WHERE id = ?`

	err := d.db.QueryRow(query, id).Scan(
		&doc.ID,
		&doc.Title,
		&doc.Description,
		&doc.FileName,
		&doc.FilePath,
		&doc.FileSize,
		&doc.FileType,
		&doc.Category,
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
	query := `SELECT id, title, COALESCE(description, '') as description, file_name, file_path,
			  file_size, file_type, COALESCE(category, '') as category, created_at, updated_at
			  FROM documents WHERE category = ? ORDER BY created_at DESC`

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

func (d *Database) Ping() error {
	return d.db.Ping()
}
