package controllers

import (
	"encoding/json"
	"net/http"

	"janus-backend-api/config"
	"janus-backend-api/middleware"
	"janus-backend-api/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ConfigController handles config management endpoints
type ConfigController struct{}

// NewConfigController creates a new ConfigController
func NewConfigController() *ConfigController {
	return &ConfigController{}
}

// List handles GET /configs - list all user's configs
func (c *ConfigController) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	var configs []models.GlobalJobConfig
	if err := config.DB.Where("user_id = ?", userID).Find(&configs).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch configs"))
		return
	}

	// Convert to response format
	responses := make([]models.ConfigResponse, len(configs))
	for i, cfg := range configs {
		responses[i] = cfg.ToResponse()
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Configs retrieved", responses))
}

// GetActive handles GET /configs/active - get currently active config
func (c *ConfigController) GetActive(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	var cfg models.GlobalJobConfig
	if err := config.DB.Where("user_id = ? AND status = ?", userID, models.ConfigStatusActive).First(&cfg).Error; err != nil {
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("No active config found"))
		return
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Active config retrieved", cfg.ToResponse()))
}

// Create handles POST /configs - create a new config
func (c *ConfigController) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	var req models.CreateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Invalid request body"))
		return
	}

	if req.ConfigName == "" {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Config name is required"))
		return
	}

	cfg := models.GlobalJobConfig{
		ConfigID:   uuid.New(),
		UserID:     userID,
		ConfigName: &req.ConfigName,
		Config:     req.Config,
		Status:     models.ConfigStatusInactive,
	}

	if err := config.DB.Create(&cfg).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to create config"))
		return
	}

	respondJSON(w, http.StatusCreated, models.NewSuccessResponse("Config created", cfg.ToResponse()))
}

// Get handles GET /configs/{id} - get config details
func (c *ConfigController) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	configID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Invalid config ID"))
		return
	}

	var cfg models.GlobalJobConfig
	if err := config.DB.Where("config_id = ? AND user_id = ?", configID, userID).First(&cfg).Error; err != nil {
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("Config not found"))
		return
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Config retrieved", cfg.ToResponse()))
}

// Update handles PUT /configs/{id} - update config
func (c *ConfigController) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	configID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Invalid config ID"))
		return
	}

	var cfg models.GlobalJobConfig
	if err := config.DB.Where("config_id = ? AND user_id = ?", configID, userID).First(&cfg).Error; err != nil {
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("Config not found"))
		return
	}

	var req models.UpdateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Invalid request body"))
		return
	}

	// Update fields
	if req.ConfigName != "" {
		cfg.ConfigName = &req.ConfigName
	}
	if req.Config != nil {
		cfg.Config = req.Config
	}

	if err := config.DB.Save(&cfg).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to update config"))
		return
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Config updated", cfg.ToResponse()))
}

// Delete handles DELETE /configs/{id} - delete config
func (c *ConfigController) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	configID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Invalid config ID"))
		return
	}

	result := config.DB.Where("config_id = ? AND user_id = ?", configID, userID).Delete(&models.GlobalJobConfig{})
	if result.Error != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to delete config"))
		return
	}
	if result.RowsAffected == 0 {
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("Config not found"))
		return
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Config deleted", nil))
}

// Activate handles POST /configs/{id}/activate - activate config
func (c *ConfigController) Activate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	configID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Invalid config ID"))
		return
	}

	// Start transaction
	tx := config.DB.Begin()

	// Deactivate all other configs for this user
	if err := tx.Model(&models.GlobalJobConfig{}).
		Where("user_id = ? AND status = ?", userID, models.ConfigStatusActive).
		Update("status", models.ConfigStatusInactive).Error; err != nil {
		tx.Rollback()
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to deactivate other configs"))
		return
	}

	// Activate the specified config
	result := tx.Model(&models.GlobalJobConfig{}).
		Where("config_id = ? AND user_id = ?", configID, userID).
		Update("status", models.ConfigStatusActive)
	if result.Error != nil {
		tx.Rollback()
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to activate config"))
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("Config not found"))
		return
	}

	tx.Commit()

	// Get updated config
	var cfg models.GlobalJobConfig
	config.DB.Where("config_id = ?", configID).First(&cfg)

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Config activated", cfg.ToResponse()))
}

// Deactivate handles POST /configs/{id}/deactivate - deactivate config
func (c *ConfigController) Deactivate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	configID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Invalid config ID"))
		return
	}

	result := config.DB.Model(&models.GlobalJobConfig{}).
		Where("config_id = ? AND user_id = ?", configID, userID).
		Update("status", models.ConfigStatusInactive)
	if result.Error != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to deactivate config"))
		return
	}
	if result.RowsAffected == 0 {
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("Config not found"))
		return
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Config deactivated", nil))
}
