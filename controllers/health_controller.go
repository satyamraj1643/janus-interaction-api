package controllers

import (
	"janus-backend-api/models"
	"net/http"
)

// HealthController handles health check endpoints
type HealthController struct{}

// NewHealthController creates a new HealthController
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Health handles GET /health
func (c *HealthController) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Status handles GET /status
func (c *HealthController) Status(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Service is running", map[string]string{
		"status":  "healthy",
		"version": "1.0.0",
	}))
}
