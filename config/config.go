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
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JanusBaseURL:       getEnv("JANUS_BASE_URL", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
	}
}
