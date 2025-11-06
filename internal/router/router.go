package router

import (
	"log"
	"net/http"

	"school-website/internal/config"
	"school-website/internal/database"
	"school-website/internal/handlers"
	"school-website/internal/middleware"
	"school-website/internal/services"

	"github.com/gorilla/mux"
)

func Setup(cfg *config.Config, db *database.Database) *mux.Router {
	r := mux.NewRouter()

	// Initialize services
	sessionService := services.NewSessionService(cfg.SessionKey)
	uploadService := services.NewFileUploadService(cfg.UploadDir)
	documentService := services.NewDocumentService(db, cfg.UploadDir+"/documents")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(sessionService.GetStore(), cfg)
	contactHandler := handlers.NewContactHandler(db)
	newsHandler := handlers.NewNewsHandler(db, uploadService)
	documentHandler := handlers.NewDocumentHandler(documentService)
	folderHandler := handlers.NewFolderHandler(db) // Добавлено

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(sessionService.GetStore())

	// --- Public Routes ---
	setupPublicRoutes(r, authHandler, contactHandler, newsHandler, documentHandler, folderHandler, cfg)

	// --- Protected Admin Routes ---
	setupAdminRoutes(r, authHandler, contactHandler, newsHandler, documentHandler, folderHandler, authMiddleware, cfg)

	// Public static files (must be last)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.PublicDir)))

	return r
}

func setupPublicRoutes(r *mux.Router, authHandler *handlers.AuthHandler,
	contactHandler *handlers.ContactHandler, newsHandler *handlers.NewsHandler,
	documentHandler *handlers.DocumentHandler, folderHandler *handlers.FolderHandler, cfg *config.Config) {

	// API endpoints
	r.HandleFunc("/api/contact", contactHandler.SubmitContact).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/news", newsHandler.GetAllNews).Methods("GET")
	r.HandleFunc("/api/news/{id}", newsHandler.GetSingleNews).Methods("GET")

	// Public document endpoints
	r.HandleFunc("/api/documents", documentHandler.GetAllDocuments).Methods("GET")
	r.HandleFunc("/api/documents/{id}", documentHandler.GetDocument).Methods("GET")
	r.HandleFunc("/api/documents/{id}/download", documentHandler.DownloadDocument).Methods("GET")

	// Public folder endpoints
	r.HandleFunc("/api/folders", folderHandler.GetAllFolders).Methods("GET")
	r.HandleFunc("/api/folders/{id}/documents", folderHandler.GetFolderDocuments).Methods("GET")

	// Auth endpoints
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Static pages
	r.HandleFunc("/admin/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.TemplatesDir+"/login.html")
	})

	r.HandleFunc("/news_article.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.PublicDir+"/news_article.html")
	})

	r.HandleFunc("/documents.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.PublicDir+"/documents.html")
	})
}

func setupAdminRoutes(r *mux.Router, authHandler *handlers.AuthHandler,
	contactHandler *handlers.ContactHandler, newsHandler *handlers.NewsHandler,
	documentHandler *handlers.DocumentHandler, folderHandler *handlers.FolderHandler,
	authMiddleware *middleware.AuthMiddleware, cfg *config.Config) {

	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(authMiddleware.RequireAuth)

	// API routes
	adminRouter.HandleFunc("/api/applications", contactHandler.GetApplications).Methods("GET", "OPTIONS")
	adminRouter.HandleFunc("/api/contacts", contactHandler.GetApplications).Methods("GET", "OPTIONS")

	// News routes
	adminRouter.HandleFunc("/api/news", newsHandler.CreateNews).Methods("POST", "OPTIONS")
	adminRouter.HandleFunc("/api/news/{id}", newsHandler.GetSingleNews).Methods("GET")
	adminRouter.HandleFunc("/api/news/{id}", newsHandler.UpdateNews).Methods("PUT")
	adminRouter.HandleFunc("/api/news/{id}", newsHandler.DeleteNews).Methods("DELETE", "OPTIONS")

	// Document routes (admin only)
	adminRouter.HandleFunc("/api/documents", documentHandler.GetAllDocuments).Methods("GET")
	adminRouter.HandleFunc("/api/documents", documentHandler.UploadDocument).Methods("POST", "OPTIONS")
	adminRouter.HandleFunc("/api/documents/{id}", documentHandler.GetDocument).Methods("GET")
	adminRouter.HandleFunc("/api/documents/{id}", documentHandler.DeleteDocument).Methods("DELETE", "OPTIONS")

	// Folder routes (admin only)
	adminRouter.HandleFunc("/api/folders", folderHandler.GetAllFolders).Methods("GET")
	adminRouter.HandleFunc("/api/folders", folderHandler.CreateFolder).Methods("POST", "OPTIONS")
	adminRouter.HandleFunc("/api/folders/{id}", folderHandler.DeleteFolder).Methods("DELETE", "OPTIONS")
	adminRouter.HandleFunc("/api/folders/{id}/documents", folderHandler.GetFolderDocuments).Methods("GET")

	// Logout
	adminRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// Static admin pages
	adminPages := map[string]string{
		"/dashboard.html":      "dashboard.html",
		"/applications.html":   "applications.html",
		"/add_news.html":       "add_news.html",
		"/news_list.html":      "news_list.html",
		"/edit_news.html":      "edit_news.html",
		"/documents_list.html": "documents_list.html",
	}

	for route, file := range adminPages {
		filePath := cfg.TemplatesDir + "/" + file
		adminRouter.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filePath)
		})
	}

	log.Println("Admin routes configured")
}
