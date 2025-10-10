package main

import (
	"log"
	"net/http"
	"os"

	"school-website/internal/config"
	"school-website/internal/database"
	"school-website/internal/router"
)

func main() {
	// Initialize directories
	initDirectories()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Setup router
	r := router.Setup(cfg, db)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Starting server on http://localhost%s", addr)
	log.Printf("Admin panel available at http://localhost%s/admin/login.html", addr)
	log.Printf("Login: %s, Password: %s", cfg.AdminUsername, cfg.AdminPassword)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDirectories() {
	dirs := []string{
		"public/uploads",
		"public/uploads/documents",
		"templates",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: failed to create directory %s: %v", dir, err)
		} else {
			log.Printf("Directory %s created or already exists", dir)
		}
	}
}
