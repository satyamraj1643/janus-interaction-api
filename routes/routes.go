package routes

import (
	"janus-backend-api/controllers"
	"janus-backend-api/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// SetupRouter configures all routes and returns the router
func SetupRouter(janusBaseURL string) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS)

	// Initialize controllers
	healthController := controllers.NewHealthController()
	authController := controllers.NewAuthController()
	submitController := controllers.NewSubmitController(janusBaseURL)
	configController := controllers.NewConfigController()
	jobController := controllers.NewJobController()
	batchController := controllers.NewBatchController()

	// ====================
	// Public Routes
	// ====================

	// Health
	r.Get("/health", healthController.Health)
	r.Get("/status", healthController.Status)

	// Auth (public)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
		r.Get("/google", authController.GoogleAuth)
		r.Get("/google/callback", authController.GoogleCallback)
	})

	// ====================
	// Protected Routes
	// ====================
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth)

		// Auth (protected)
		r.Get("/auth/profile", authController.Profile)

		// Job Submission (Proxy to Janus)
		r.Route("/submit", func(r chi.Router) {
			r.Post("/job", submitController.SubmitJob)
			r.Post("/batch", submitController.SubmitBatch)
			r.Post("/batch/atomic", submitController.SubmitBatchAtomic)
		})

		// Config Management
		r.Route("/configs", func(r chi.Router) {
			r.Get("/", configController.List)
			r.Get("/active", configController.GetActive)
			r.Post("/", configController.Create)
			r.Get("/{id}", configController.Get)
			r.Put("/{id}", configController.Update)
			r.Delete("/{id}", configController.Delete)
			r.Post("/{id}/activate", configController.Activate)
			r.Post("/{id}/deactivate", configController.Deactivate)
		})

		// Jobs
		r.Route("/jobs", func(r chi.Router) {
			r.Get("/", jobController.List)
			r.Get("/stats", jobController.Stats)
			r.Get("/{id}", jobController.Get)
		})

		// Batches
		r.Route("/batches", func(r chi.Router) {
			r.Get("/", batchController.List)
			r.Get("/{id}", batchController.Get)
			r.Get("/{id}/jobs", batchController.GetJobs)
		})
	})

	return r
}
