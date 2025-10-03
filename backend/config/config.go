package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Port          string
	Env           string
	IP            string
	Database      DatabaseConfig
	RAGServiceURL string
	JWT           JWTConfig
	Upload        UploadConfig
	CORS          CORSConfig
	RateLimit     RateLimitConfig
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	MaxSize      int64
	AllowedTypes string
}

// CORSConfig CORS 配置
type CORSConfig struct {
	Origins string
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Requests int
	Window   int
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	Charset     string
	PoolSize    int
	MaxOverflow int
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 首先尝试加载外层的 .env 文件
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Could not load outer .env file: %v", err)
	}

	// 然后加载当前目录的 .env 文件（会覆盖外层 .env 中的同名变量）
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load local .env file: %v", err)

	}

	return &Config{
		Port: getEnv("API_PORT", "9000"),
		IP:   getEnv("IP", "127.0.0.1"),
		Env:  getEnv("BACKEND_ENV", "development"),
		Database: DatabaseConfig{
			Host:        getEnv("DATABASE_HOST", "127.0.0.1"),
			Port:        getEnv("DATABASE_PORT", "3306"),
			User:        getEnv("DATABASE_USERNAME", "root"),
			Password:    getEnv("DATABASE_PASSWORD", "123456"),
			DBName:      getEnv("DATABASE_NAME", "xxx_review"),
			Charset:     getEnv("DATABASE_CHARSET", "utf8mb4"),
			PoolSize:    getEnvAsInt("DATABASE_POOL_SIZE", 10),
			MaxOverflow: getEnvAsInt("DATABASE_MAX_OVERFLOW", 20),
		},
		RAGServiceURL: getEnv("RAG_SERVICE_URL", "http://localhost:50051"),
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your_jwt_secret_key_here"),
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		},
		Upload: UploadConfig{
			MaxSize:      int64(getEnvAsInt("UPLOAD_MAX_SIZE", 10485760)),
			AllowedTypes: getEnv("UPLOAD_ALLOWED_TYPES", "pdf,doc,docx,txt,md"),
		},
		CORS: CORSConfig{
			Origins: getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:9000"),
		},
		RateLimit: RateLimitConfig{
			Requests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			Window:   getEnvAsInt("RATE_LIMIT_WINDOW", 60),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt 获取环境变量并转换为整数，如果不存在或转换失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}
