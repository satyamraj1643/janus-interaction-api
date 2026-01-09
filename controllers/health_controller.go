package controllers

import (
	"janus-backend-api/models"
	"net/http"
)

// HealthController handles health check endpoints
type HealthController struct {
	janusBaseURL string
}

// NewHealthController creates a new HealthController
func NewHealthController(janusBaseURL string) *HealthController {
	return &HealthController{
		janusBaseURL: janusBaseURL,
	}
}

// Health handles GET /health
func (c *HealthController) Health(w http.ResponseWriter, r *http.Request) {
	apiStatus := "ok"
	janusStatus := "unknown"

	// Check Janus Health
	resp, err := http.Get(c.janusBaseURL + "/health")
	if err == nil {
		if resp.StatusCode == http.StatusOK {
			janusStatus = "ok"
		} else {
			janusStatus = "unhealthy"
		}
		resp.Body.Close()
	} else {
		janusStatus = "down"
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"api_status":   apiStatus,
		"janus_status": janusStatus,
	})
}

// Status handles GET /status
func (c *HealthController) Status(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Service is running", map[string]string{
		"status":  "healthy",
		"version": "1.0.0",
	}))
}
