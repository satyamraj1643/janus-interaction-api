package controllers

import (
	"net/http"
	"strconv"

	"janus-backend-api/config"
	"janus-backend-api/middleware"
	"janus-backend-api/models"

	"github.com/go-chi/chi/v5"
)

// JobController handles job viewing endpoints
type JobController struct{}

// NewJobController creates a new JobController
func NewJobController() *JobController {
	return &JobController{}
}

// List handles GET /jobs - list all user's jobs with pagination
func (c *JobController) List(w http.ResponseWriter, r *http.Request) {
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

	// Optional filters
	status := r.URL.Query().Get("status")
	batchID := r.URL.Query().Get("batch_id")

	// Build query
	query := config.DB.Model(&models.Job{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("job_status = ?", status)
	}
	if batchID != "" {
		query = query.Where("batch_id = ?", batchID)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Fetch jobs
	var jobs []models.Job
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&jobs).Error; err != nil {
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

// Get handles GET /jobs/{id} - get job details
func (c *JobController) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	jobID := chi.URLParam(r, "id")

	var job models.Job
	if err := config.DB.Where("job_id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("Job not found"))
		return
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Job retrieved", job.ToResponse()))
}

// Stats handles GET /jobs/stats - get job statistics
func (c *JobController) Stats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	var stats models.StatsResponse

	// Total jobs
	config.DB.Model(&models.Job{}).Where("user_id = ?", userID).Count(&stats.TotalJobs)

	// Accepted jobs
	config.DB.Model(&models.Job{}).Where("user_id = ? AND job_status = ?", userID, "accepted").Count(&stats.AcceptedJobs)

	// Rejected jobs
	config.DB.Model(&models.Job{}).Where("user_id = ? AND job_status = ?", userID, "rejected").Count(&stats.RejectedJobs)

	// Total batches
	config.DB.Model(&models.Batch{}).Where("user_id = ?", userID).Count(&stats.TotalBatches)

	// Total configs
	config.DB.Model(&models.GlobalJobConfig{}).Where("user_id = ?", userID).Count(&stats.TotalConfigs)

	// Active configs
	config.DB.Model(&models.GlobalJobConfig{}).Where("user_id = ? AND status = ?", userID, models.ConfigStatusActive).Count(&stats.ActiveConfigs)

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Stats retrieved", stats))
}
