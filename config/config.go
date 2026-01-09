package config

// AppConfig holds application configuration
type AppConfig struct {
	ServerPort         string
	DatabaseURL        string
	JWTSecret          string
	JanusBaseURL       string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *AppConfig {
	return &AppConfig{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgresql://janus_admin:t6exln67eDs6ygXQQTkaJE2CSsBvtMtl@dpg-d571ksuuk2gs73cp2sjg-a.oregon-postgres.render.com/janus_db_03vr?sslmode=require"),
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JanusBaseURL:       getEnv("JANUS_BASE_URL", "https://janus-microservice.onrender.com"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
	}
}
