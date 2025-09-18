package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ContactForm определяет структуру данных, получаемых из формы
type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS для запросов с фронтенда (на время разработки)
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var form ContactForm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&form)
	if err != nil {
		http.Error(w, "Ошибка при чтении данных", http.StatusBadRequest)
		return
	}

	// Сохраняем данные в базу данных
	err = SaveContact(form)
	if err != nil {
		http.Error(w, "Ошибка при сохранении данных в базу", http.StatusInternalServerError)
		return
	}

	log.Printf("Сообщение от %s (%s) успешно сохранено в БД.", form.Name, form.Email)

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Сообщение успешно получено и сохранено!"})
}

func main() {
	// Инициализируем базу данных
	InitDB("school.db")

	// Создаем обработчик для статических файлов (наш index.html)
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// Регистрируем обработчик для API
	http.HandleFunc("/api/contact", contactHandler)

	// Запускаем сервер на порту 8080
	log.Println("Запуск сервера на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
