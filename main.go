package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
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
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// Credentials для парсинга JSON при логине
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func handleFileUpload(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil // Файл не был загружен, это не ошибка
		}
		return "", err
	}
	defer file.Close()

	// Создаем уникальное имя файла
	uniqueFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
	filePath := filepath.Join("public", "uploads", uniqueFileName)

	// Создаем файл на сервере
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Копируем содержимое загруженного файла
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/uploads/" + uniqueFileName, nil
}

func adminCreateNewsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Получаем текстовые поля из формы
	title := r.FormValue("title")
	content := r.FormValue("content")
	imageURLFromForm := r.FormValue("image_url")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Приоритет у загруженного файла
	finalImageURL, err := handleFileUpload(r)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	if finalImageURL == "" {
		finalImageURL = imageURLFromForm // Если файл не загружен, используем URL из формы
	}

	err = SaveNews(title, content, finalImageURL)
	if err != nil {
		http.Error(w, "Failed to save news", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getSingleNewsArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing news ID", http.StatusBadRequest)
		return
	}

	article, err := GetNewsArticle(id)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func updateNewsArticleHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing news ID", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Получаем текущую статью из БД, чтобы знать ее текущий URL изображения
	existingArticle, err := GetNewsArticle(id)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// Пытаемся обработать загрузку нового файла
	newImageURL, err := handleFileUpload(r)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	// Определяем финальный URL изображения
	finalImageURL := existingArticle.ImageURL // 1. По умолчанию, сохраняем старый URL

	if newImageURL != "" {
		// 2. Если загружен новый файл, он имеет наивысший приоритет
		finalImageURL = newImageURL
	} else if imageURLFromForm != existingArticle.ImageURL {
		// 3. Если файл не загружен, проверяем, изменил ли пользователь текстовое поле URL.
		// Это позволяет пользователю как обновить URL на новый, так и удалить его (отправив пустое поле),
		// но не удалит существующий URL, если пользователь просто ничего не трогал.
		finalImageURL = imageURLFromForm
	}

	// Обновляем статью в БД
	err = UpdateNewsArticle(id, title, content, finalImageURL)
	if err != nil {
		http.Error(w, "Failed to update news", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	r := mux.NewRouter()

	// --- Публичные API и страницы ---
	r.HandleFunc("/api/contact", contactHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/news", newsAPIHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/admin/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/login.html")
	})

	// --- Защищенный Subrouter для админ-панели ---
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(authMiddleware) // Применяем middleware ко всем маршрутам админки

	// API маршруты админки
	adminRouter.HandleFunc("/api/applications", applicationsAPIHandler).Methods("GET")
	adminRouter.HandleFunc("/api/news", adminCreateNewsHandler).Methods("POST")
	adminRouter.HandleFunc("/api/news/{id}", getSingleNewsArticleHandler).Methods("GET")
	adminRouter.HandleFunc("/api/news/{id}", updateNewsArticleHandler).Methods("PUT")
	adminRouter.HandleFunc("/logout", logoutHandler).Methods("POST")

	// Статические файлы админки
	adminRouter.HandleFunc("/dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/dashboard.html")
	})
	adminRouter.HandleFunc("/applications.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/applications.html")
	})
	adminRouter.HandleFunc("/add_news.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/add_news.html")
	})
	adminRouter.HandleFunc("/news_list.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/news_list.html")
	})
	adminRouter.HandleFunc("/edit_news.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/edit_news.html")
	})


	// --- Публичный статический сайт ---
	// Этот обработчик должен быть последним
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	// Запускаем сервер на порту 8080
	log.Println("Запуск сервера на http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
