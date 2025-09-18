package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // Импортируем драйвер SQLite
)

var db *sql.DB // Глобальная переменная для подключения к БД

// InitDB инициализирует подключение к базе данных и создает таблицу, если ее нет.
func InitDB(filepath string) {
	var err error
	// Открываем или создаем файл БД
	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Ошибка при открытии БД: %v", err)
	}

	// SQL-запрос для создания таблицы
	createTableSQL := `CREATE TABLE IF NOT EXISTS contacts (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"email" TEXT,
		"message" TEXT,
		"created_at" DATETIME
	);`

	// Выполняем SQL-запрос
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatalf("Ошибка при подготовке SQL-запроса для создания таблицы: %v", err)
	}
	statement.Exec()
	log.Println("Таблица 'contacts' успешно создана или уже существует.")
}

// ContactEntry представляет одну запись в таблице contacts
type ContactEntry struct {
	ID        int
	Name      string
	Email     string
	Message   string
	CreatedAt time.Time
}

// SaveContact сохраняет данные из формы в базу данных.
func SaveContact(form ContactForm) error {
	insertSQL := `INSERT INTO contacts(name, email, message, created_at) VALUES (?, ?, ?, ?)`

	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Printf("Ошибка при подготовке SQL-запроса для вставки: %v", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(form.Name, form.Email, form.Message, time.Now())
	if err != nil {
		log.Printf("Ошибка при выполнении SQL-запроса для вставки: %v", err)
		return err
	}

	return nil
}

// GetContacts извлекает все заявки из базы данных.
func GetContacts() ([]ContactEntry, error) {
	rows, err := db.Query("SELECT id, name, email, message, created_at FROM contacts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []ContactEntry
	for rows.Next() {
		var c ContactEntry
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Message, &c.CreatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}

	return contacts, nil
}
