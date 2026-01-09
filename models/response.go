package models

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(err string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   err,
	}
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, page, perPage int, totalItems int64) PaginatedResponse {
	totalPages := int(totalItems) / perPage
	if int(totalItems)%perPage > 0 {
		totalPages++
	}
	return PaginatedResponse{
		Success:    true,
		Data:       data,
		Page:       page,
		PerPage:    perPage,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

// StatsResponse for job/batch statistics
type StatsResponse struct {
	TotalJobs     int64 `json:"total_jobs"`
	AcceptedJobs  int64 `json:"accepted_jobs"`
	RejectedJobs  int64 `json:"rejected_jobs"`
	TotalBatches  int64 `json:"total_batches"`
	TotalConfigs  int64 `json:"total_configs"`
	ActiveConfigs int64 `json:"active_configs"`
}
