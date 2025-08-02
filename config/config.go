package config

import (
	"os"
	"time"
)

// Config holds all configuration for our application

type SendGridConfig struct {
	APIKey string
}

type Config struct {
	Server    ServerConfig
	App       AppConfig
	Database  DatabaseConfig
	SendGrid  SendGridConfig
	SecretKey SecretKeyConfig
	Line      LineConfig
}

// ServerConfig holds all the server-related config
type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	AppName        string
	AllowedOrigins string
}

// AppConfig holds all the application-related config
type AppConfig struct {
	Environment      string
	AppName          string
	FrontendURLS     string
	UploadPath       string
	APIBaseURL       string
	StaticFilePrefix string
	CMSBaseURL       string `json:"cms_base_url" env:"CMS_BASE_URL"`
	WebBaseURL       string `json:"web_base_url" env:"WEB_BASE_URL"`
}

// DatabaseConfig holds all the database-related config
type DatabaseConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

type SecretKeyConfig struct {
	LineKey   string
	NormalKey string
}

type LineConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUri  string
	TokenUrl     string
	AuthorizeUrl string
}

func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			AppName: getEnv("APP_NAME", "CMS API"),
		},
		App: AppConfig{
			Environment:      getEnv("ENVIRONMENT", "development"),
			FrontendURLS:     getEnv("FRONTEND_URLS", "http://localhost:3000,http://localhost:3001"),
			AppName:          getEnv("APP_NAME_APP", "CMS Application"),
			UploadPath:       getEnv("UPLOAD_FILE_PATH", "./public_uploads"),
			APIBaseURL:       getEnv("API_BASE_URL", "http://localhost:8080"),
			StaticFilePrefix: getEnv("STATIC_FILE_PREFIX", "/files"),
			CMSBaseURL:       getEnv("CMS_BASE_URL", "http://localhost:3001"),
			WebBaseURL:       getEnv("WEB_BASE_URL", "http://localhost:3030"),
		},
		Database: DatabaseConfig{

			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			Username:     getEnv("DB_USERNAME", ""),
			Password:     getEnv("DB_PASSWORD", ""),
			DatabaseName: getEnv("DB_NAME", ""),
		},
		SendGrid: SendGridConfig{

			APIKey: getEnv("SENDGRID_API_KEY", ""),
		},
		SecretKey: SecretKeyConfig{
			LineKey:   getEnv("OAUTH_CLIENT_SECRET", ""),
			NormalKey: getEnv("JWT_SECRET_KEY", ""),
		},
		Line: LineConfig{
			ClientId:     getEnv("OAUTH_CLIENT_ID", ""),
			ClientSecret: getEnv("OAUTH_CLIENT_SECRET", ""),
			RedirectUri:  getEnv("OAUTH_REDIRECT_URI", "http://localhost:8080/auth/callback"),
			TokenUrl:     getEnv("TOKEN_URL", "https://api.line.biz/oauth2/v2.1/token"),
			AuthorizeUrl: getEnv("AUTHORIZE_URL", "https://access.line.me/oauth2/v2.1/authorize"),
		},
	}
}

// Simple helper function to read an environment variable or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
