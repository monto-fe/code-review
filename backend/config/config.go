package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Port     string
	Env      string
	IP       string
	Database DatabaseConfig
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 首先尝试加载外层的 .env 文件
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: Could not load outer .env file: %v", err)
	}

	// 然后加载当前目录的 .env 文件（会覆盖外层 .env 中的同名变量）
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load local .env file: %v", err)
	}

	return &Config{
		Port: getEnv("API_PORT", "9000"),
		IP:   getEnv("IP", "127.0.0.1"),
		Env:  getEnv("NODE_ENV", "development"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_DATABASE", "code_review"),
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
