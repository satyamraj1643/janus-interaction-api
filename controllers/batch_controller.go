package controllers

import (
	"net/http"
	"strconv"

	"janus-backend-api/config"
	"janus-backend-api/middleware"
	"janus-backend-api/models"

	"github.com/go-chi/chi/v5"
)

// BatchController handles batch viewing endpoints
type BatchController struct{}

// NewBatchController creates a new BatchController
func NewBatchController() *BatchController {
	return &BatchController{}
}

// List handles GET /batches - list all user's batches
func (c *BatchController) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	// Parse pagination params
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Count total
	var total int64
	config.DB.Model(&models.Batch{}).Where("user_id = ?", userID).Count(&total)

	// Fetch batches
	var batches []models.Batch
	offset := (page - 1) * perPage
	if err := config.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&batches).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch batches"))
		return
	}

	// Convert to response format
	responses := make([]models.BatchResponse, len(batches))
	for i, batch := range batches {
		responses[i] = batch.ToResponse()
	}

	respondJSON(w, http.StatusOK, models.NewPaginatedResponse(responses, page, perPage, total))
}

// Get handles GET /batches/{id} - get batch with summary
func (c *BatchController) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	batchID := chi.URLParam(r, "id")

	var batch models.Batch
	if err := config.DB.Where("batch_id = ? AND user_id = ?", batchID, userID).First(&batch).Error; err != nil {
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("Batch not found"))
		return
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Batch retrieved", batch.ToResponse()))
}

// GetJobs handles GET /batches/{id}/jobs - get jobs in a batch
func (c *BatchController) GetJobs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	batchID := chi.URLParam(r, "id")

	// Verify batch belongs to user
	var batch models.Batch
	if err := config.DB.Where("batch_id = ? AND user_id = ?", batchID, userID).First(&batch).Error; err != nil {
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("Batch not found"))
		return
	}

	// Parse pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}

	// Count total
	var total int64
	config.DB.Model(&models.Job{}).Where("batch_id = ?", batchID).Count(&total)

	// Fetch jobs
	var jobs []models.Job
	offset := (page - 1) * perPage
	if err := config.DB.Where("batch_id = ?", batchID).
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&jobs).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch jobs"))
		return
	}

	// Convert to response format
	responses := make([]models.JobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = job.ToResponse()
	}

	respondJSON(w, http.StatusOK, models.NewPaginatedResponse(responses, page, perPage, total))
}
