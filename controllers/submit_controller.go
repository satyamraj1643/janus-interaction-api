package controllers

import (
	"bytes"
	"io"
	"net/http"

	"janus-backend-api/middleware"
	"janus-backend-api/models"
)

// SubmitController handles job submission proxy to Janus
type SubmitController struct {
	janusBaseURL string
	httpClient   *http.Client
}

// NewSubmitController creates a new SubmitController
func NewSubmitController(janusBaseURL string) *SubmitController {
	return &SubmitController{
		janusBaseURL: janusBaseURL,
		httpClient:   &http.Client{},
	}
}

// SubmitJob handles POST /submit/job - proxies to Janus /dashboard/jobs
func (c *SubmitController) SubmitJob(w http.ResponseWriter, r *http.Request) {
	c.proxyToJanus(w, r, "/dashboard/jobs")
}

// SubmitBatch handles POST /submit/batch - proxies to Janus /dashboard/jobs/batch
func (c *SubmitController) SubmitBatch(w http.ResponseWriter, r *http.Request) {
	c.proxyToJanus(w, r, "/dashboard/jobs/batch")
}

// SubmitBatchAtomic handles POST /submit/batch/atomic - proxies to Janus /dashboard/jobs/batch/atomic
func (c *SubmitController) SubmitBatchAtomic(w http.ResponseWriter, r *http.Request) {
	c.proxyToJanus(w, r, "/dashboard/jobs/batch/atomic")
}

// proxyToJanus forwards the request to the Janus microservice
func (c *SubmitController) proxyToJanus(w http.ResponseWriter, r *http.Request, path string) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Failed to read request body"))
		return
	}
	defer r.Body.Close()

	// Create proxy request
	proxyURL := c.janusBaseURL + path
	proxyReq, err := http.NewRequest(http.MethodPost, proxyURL, bytes.NewReader(body))
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to create proxy request"))
		return
	}

	// Set headers
	proxyReq.Header.Set("Content-Type", "application/json")
	proxyReq.Header.Set("X-User-ID", userID.String())

	// Forward the request
	resp, err := c.httpClient.Do(proxyReq)
	if err != nil {
		respondJSON(w, http.StatusBadGateway, models.NewErrorResponse("Failed to connect to Janus service: "+err.Error()))
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to read Janus response"))
		return
	}

	// Forward the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

// SubmitJobRequest represents a job submission request
type SubmitJobRequest struct {
	BatchName    string                 `json:"batch_name"`
	TenantID     string                 `json:"tenant_id"`
	Priority     int                    `json:"priority"`
	Dependencies map[string]int         `json:"dependencies,omitempty"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
}

// SubmitBatchRequest represents a batch submission request
type SubmitBatchRequest struct {
	BatchName string         `json:"batch_name"`
	Jobs      []BatchJobItem `json:"jobs"`
}

// BatchJobItem represents a single job in a batch
type BatchJobItem struct {
	TenantID     string                 `json:"tenant_id"`
	Priority     int                    `json:"priority"`
	Dependencies map[string]int         `json:"dependencies,omitempty"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
}
