package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ConfigStatus represents the status of a config
type ConfigStatus string

const (
	ConfigStatusActive   ConfigStatus = "active"
	ConfigStatusInactive ConfigStatus = "inactive"
)

// JSONB type for JSON columns
type JSONB map[string]interface{}

// Value implements driver.Valuer for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONB")
	}
	return json.Unmarshal(bytes, j)
}

// GlobalJobConfig represents a job configuration
type GlobalJobConfig struct {
	ConfigID   uuid.UUID    `json:"config_id" gorm:"type:uuid;primaryKey;column:config_id"`
	UserID     uuid.UUID    `json:"user_id" gorm:"type:uuid;column:user_id"`
	ConfigName *string      `json:"config_name" gorm:"column:config_name"`
	Config     JSONB        `json:"config" gorm:"type:json;column:config"`
	Status     ConfigStatus `json:"status" gorm:"type:config_status;column:status"`
}

// TableName specifies the table name for GORM
func (GlobalJobConfig) TableName() string {
	return "global_job_config"
}

// CreateConfigRequest for creating a new config
type CreateConfigRequest struct {
	ConfigName string                 `json:"config_name" binding:"required"`
	Config     map[string]interface{} `json:"config" binding:"required"`
}

// UpdateConfigRequest for updating an existing config
type UpdateConfigRequest struct {
	ConfigName string                 `json:"config_name,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

// ConfigResponse returned to clients
type ConfigResponse struct {
	ConfigID   uuid.UUID              `json:"config_id"`
	ConfigName string                 `json:"config_name"`
	Config     map[string]interface{} `json:"config"`
	Status     string                 `json:"status"`
	IsActive   bool                   `json:"is_active"`
}

// ToResponse converts GlobalJobConfig to ConfigResponse
func (c *GlobalJobConfig) ToResponse() ConfigResponse {
	name := ""
	if c.ConfigName != nil {
		name = *c.ConfigName
	}
	return ConfigResponse{
		ConfigID:   c.ConfigID,
		ConfigName: name,
		Config:     c.Config,
		Status:     string(c.Status),
		IsActive:   c.Status == ConfigStatusActive,
	}
}

// ServiceStatus represents the service status for a user
type ServiceStatus struct {
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey;column:user_id"`
	Status string    `json:"status" gorm:"column:status"`
}

// TableName specifies the table name for GORM
func (ServiceStatus) TableName() string {
	return "service_status"
}

// UserAssociation represents user statistics per config
type UserAssociation struct {
	ConfigID      uuid.UUID `json:"config_id" gorm:"type:uuid;column:config_id"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
	NoOfBatches   *int      `json:"no_of_batches" gorm:"column:no_of_batches"`
	NoOfJobs      *int      `json:"no_of_jobs" gorm:"column:no_of_jobs"`
	SucceededJobs *int      `json:"succeeded_jobs" gorm:"column:succeeded_jobs"`
	FailedJobs    *int      `json:"failed_jobs" gorm:"column:failed_jobs"`
	TotalJobs     *int      `json:"total_jobs" gorm:"column:total_jobs"`
}

// TableName specifies the table name for GORM
func (UserAssociation) TableName() string {
	return "user_association"
}

// Job represents a job in the system
type Job struct {
	JobID          string     `json:"job_id" gorm:"type:text;primaryKey;column:job_id"`
	UserID         uuid.UUID  `json:"user_id" gorm:"type:uuid;column:user_id"`
	JobPayload     JSONB      `json:"job_payload" gorm:"type:json;column:job_payload"`
	BatchID        *string    `json:"batch_id" gorm:"column:batch_id"`
	JobStatus      string     `json:"job_status" gorm:"type:job_status;column:job_status"`
	Reason         *string    `json:"reason" gorm:"column:reason"`
	CreatedAt      *time.Time `json:"created_at" gorm:"column:created_at"`
	GlobalConfigID *uuid.UUID `json:"global_config_id" gorm:"type:uuid;column:global_config_id"`
}

// TableName specifies the table name for GORM
func (Job) TableName() string {
	return "jobs"
}

// JobResponse returned to clients
type JobResponse struct {
	JobID          string                 `json:"job_id"`
	UserID         uuid.UUID              `json:"user_id"`
	JobPayload     map[string]interface{} `json:"job_payload"`
	BatchID        string                 `json:"batch_id,omitempty"`
	JobStatus      string                 `json:"job_status"`
	Reason         string                 `json:"reason,omitempty"`
	CreatedAt      *time.Time             `json:"created_at"`
	GlobalConfigID string                 `json:"global_config_id,omitempty"`
}

// ToResponse converts Job to JobResponse
func (j *Job) ToResponse() JobResponse {
	batchID := ""
	if j.BatchID != nil {
		batchID = *j.BatchID
	}
	reason := ""
	if j.Reason != nil {
		reason = *j.Reason
	}
	configID := ""
	if j.GlobalConfigID != nil {
		configID = j.GlobalConfigID.String()
	}
	return JobResponse{
		JobID:          j.JobID,
		UserID:         j.UserID,
		JobPayload:     j.JobPayload,
		BatchID:        batchID,
		JobStatus:      j.JobStatus,
		Reason:         reason,
		CreatedAt:      j.CreatedAt,
		GlobalConfigID: configID,
	}
}

// Batch represents a batch of jobs
type Batch struct {
	BatchID      string     `json:"batch_id" gorm:"type:text;primaryKey;column:batch_id"`
	BatchName    *string    `json:"batch_name" gorm:"column:batch_name"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;column:user_id"`
	CreatedAt    *time.Time `json:"created_at" gorm:"column:created_at"`
	TotalJobs    *int       `json:"total_jobs" gorm:"column:total_jobs"`
	AdmittedJobs *int       `json:"admitted_jobs" gorm:"column:admitted_jobs"`
}

// TableName specifies the table name for GORM
func (Batch) TableName() string {
	return "batch"
}

// BatchResponse returned to clients
type BatchResponse struct {
	BatchID      string     `json:"batch_id"`
	BatchName    string     `json:"batch_name"`
	UserID       uuid.UUID  `json:"user_id"`
	CreatedAt    *time.Time `json:"created_at"`
	TotalJobs    int        `json:"total_jobs"`
	AdmittedJobs int        `json:"admitted_jobs"`
	RejectedJobs int        `json:"rejected_jobs"`
}

// ToResponse converts Batch to BatchResponse
func (b *Batch) ToResponse() BatchResponse {
	name := ""
	if b.BatchName != nil {
		name = *b.BatchName
	}
	total := 0
	if b.TotalJobs != nil {
		total = *b.TotalJobs
	}
	admitted := 0
	if b.AdmittedJobs != nil {
		admitted = *b.AdmittedJobs
	}
	return BatchResponse{
		BatchID:      b.BatchID,
		BatchName:    name,
		UserID:       b.UserID,
		CreatedAt:    b.CreatedAt,
		TotalJobs:    total,
		AdmittedJobs: admitted,
		RejectedJobs: total - admitted,
	}
}
