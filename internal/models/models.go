package models

import "time"

// ContactForm defines the structure of data received from the form
type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// ContactEntry represents a single record in the contacts table
type ContactEntry struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// NewsArticle represents a single news article
type NewsArticle struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

// Credentials for parsing JSON during login
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
