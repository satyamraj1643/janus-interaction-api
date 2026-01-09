package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;primaryKey;column:user_id"`
	Name         string     `json:"name" gorm:"column:name"`
	Email        *string    `json:"email,omitempty" gorm:"column:email;uniqueIndex"`
	PasswordHash *string    `json:"-" gorm:"column:password_hash"`
	GoogleID     *string    `json:"google_id,omitempty" gorm:"column:google_id;uniqueIndex"`
	CreatedAt    *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// UserResponse is the safe user data returned to clients
type UserResponse struct {
	UserID    uuid.UUID  `json:"user_id"`
	Name      string     `json:"name"`
	Email     string     `json:"email,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// ToResponse converts User to UserResponse (hides sensitive fields)
func (u *User) ToResponse() UserResponse {
	email := ""
	if u.Email != nil {
		email = *u.Email
	}
	return UserResponse{
		UserID:    u.UserID,
		Name:      u.Name,
		Email:     email,
		CreatedAt: u.CreatedAt,
	}
}

// RegisterRequest for user registration
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse returned after successful auth
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
