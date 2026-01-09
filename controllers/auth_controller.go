package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"janus-backend-api/config"
	"janus-backend-api/middleware"
	"janus-backend-api/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthController handles authentication endpoints
type AuthController struct{}

// NewAuthController creates a new AuthController
func NewAuthController() *AuthController {
	return &AuthController{}
}

// Register handles POST /auth/register
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Invalid request body"))
		return
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Password == "" {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Name, email, and password are required"))
		return
	}

	if len(req.Password) < 6 {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Password must be at least 6 characters"))
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		respondJSON(w, http.StatusConflict, models.NewErrorResponse("Email already registered"))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to process password"))
		return
	}

	// Create user
	now := time.Now()
	passwordHash := string(hashedPassword)
	user := models.User{
		UserID:       uuid.New(),
		Name:         req.Name,
		Email:        &req.Email,
		PasswordHash: &passwordHash,
		CreatedAt:    &now,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to create user"))
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.UserID, req.Email)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to generate token"))
		return
	}

	respondJSON(w, http.StatusCreated, models.NewSuccessResponse("User registered successfully", models.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}))
}

// Login handles POST /auth/login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Invalid request body"))
		return
	}

	if req.Email == "" || req.Password == "" {
		respondJSON(w, http.StatusBadRequest, models.NewErrorResponse("Email and password are required"))
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("Invalid email or password"))
		return
	}

	// Check password
	if user.PasswordHash == nil {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("This account uses Google login"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("Invalid email or password"))
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.UserID, req.Email)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, models.NewErrorResponse("Failed to generate token"))
		return
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Login successful", models.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}))
}

// Profile handles GET /auth/profile
func (c *AuthController) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, models.NewErrorResponse("User not authenticated"))
		return
	}

	var user models.User
	if err := config.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		respondJSON(w, http.StatusNotFound, models.NewErrorResponse("User not found"))
		return
	}

	respondJSON(w, http.StatusOK, models.NewSuccessResponse("Profile retrieved", user.ToResponse()))
}

// GoogleAuth handles GET /auth/google (placeholder)
func (c *AuthController) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Google OAuth redirect
	respondJSON(w, http.StatusNotImplemented, models.NewErrorResponse("Google OAuth not configured. Please set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET"))
}

// GoogleCallback handles GET /auth/google/callback (placeholder)
func (c *AuthController) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Google OAuth callback
	respondJSON(w, http.StatusNotImplemented, models.NewErrorResponse("Google OAuth not configured"))
}
