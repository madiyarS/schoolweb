package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

// --- Глобальные переменные ---
var (
	// Ключ для cookie. ВАЖНО: В реальном приложении его нужно генерировать (например, `openssl rand -hex 32`) и хранить безопасно, например, в переменных окружения.
	key   = []byte("super-secret-key-that-is-32-bytes-long-so-it-is-secure")
	store = sessions.NewCookieStore(key)

	// Учетные данные администратора (для простоты жестко закодированы)
	adminUsername = "admin"
	adminPassword = "password123"
)

// --- Структуры данных ---

// ContactForm определяет структуру данных, получаемых из формы
type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// Credentials для парсинга JSON при логине
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewsFormData для парсинга JSON при создании новости
type NewsFormData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}


// --- Middleware ---

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		// Проверяем, аутентифицирован ли пользователь
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// Перенаправляем на страницу входа, если не аутентифицирован
			http.Redirect(w, r, "/admin/login.html", http.StatusFound)
			return
		}

		// Если да, передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}


// --- Обработчики HTTP ---

func newsAPIHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := GetNews()
	if err != nil {
		http.Error(w, "Failed to get news", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func adminNewsAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var formData NewsFormData
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if formData.Title == "" || formData.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	err = SaveNews(formData.Title, formData.Content)
	if err != nil {
		http.Error(w, "Failed to save news", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func applicationsAPIHandler(w http.ResponseWriter, r *http.Request) {
	contacts, err := GetContacts()
	if err != nil {
		http.Error(w, "Failed to get contacts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Проверяем логин и пароль
	if creds.Username != adminUsername || creds.Password != adminPassword {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Создаем сессию
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	// Сбрасываем флаг аутентификации
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1 // Удаляем cookie

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/login.html", http.StatusFound)
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

	// --- Настройка маршрутизатора ---
	mux := http.NewServeMux()

	// --- Публичные маршруты ---
	mux.HandleFunc("/api/contact", contactHandler)
	mux.HandleFunc("/api/news", newsAPIHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/admin/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/login.html")
	})


	// --- Защищенные маршруты админ-панели ---
	// Обратите внимание: мы больше не используем отдельный adminMux, а регистрируем все на основном mux
	// и оборачиваем каждый защищенный маршрут в authMiddleware.
	mux.Handle("/admin/logout", authMiddleware(http.HandlerFunc(logoutHandler)))
	mux.Handle("/admin/api/applications", authMiddleware(http.HandlerFunc(applicationsAPIHandler)))
	mux.Handle("/admin/api/news", authMiddleware(http.HandlerFunc(adminNewsAPIHandler)))

	mux.Handle("/admin/dashboard.html", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/dashboard.html")
	})))
	mux.Handle("/admin/applications.html", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/applications.html")
	})))
	mux.Handle("/admin/add_news.html", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/add_news.html")
	})))


	// --- Публичный статический сайт ---
	// Этот обработчик должен быть последним, так как он ловит все остальные запросы
	mux.Handle("/", http.FileServer(http.Dir("public")))


	// Запускаем сервер на порту 8080
	log.Println("Запуск сервера на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
