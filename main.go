package main

import (
	"fmt"
	"log"
	"net/http"

	"janus-backend-api/config"
	"janus-backend-api/middleware"
	"janus-backend-api/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set JWT secret
	middleware.SetJWTSecret(cfg.JWTSecret)

	// Connect to database
	config.ConnectDatabase()

	// Run migrations for auth columns (if they don't exist)
	runMigrations()

	// Setup router
	router := routes.SetupRouter(cfg.JanusBaseURL)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Janus API starting on http://localhost%s", addr)
	log.Printf("📍 API Endpoints:")
	log.Printf("   Auth:    /auth/register, /auth/login, /auth/profile")
	log.Printf("   Submit:  /submit/job, /submit/batch, /submit/batch/atomic")
	log.Printf("   Configs: /configs (CRUD + activate/deactivate)")
	log.Printf("   Jobs:    /jobs, /jobs/stats, /jobs/{id}")
	log.Printf("   Batches: /batches, /batches/{id}, /batches/{id}/jobs")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// runMigrations adds auth columns to users table if they don't exist
func runMigrations() {
	// Add email column if not exists
	config.DB.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT UNIQUE;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id TEXT UNIQUE;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();
	`)
	log.Println("✅ Database migrations complete")
}
