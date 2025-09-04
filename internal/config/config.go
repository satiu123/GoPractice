package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Redis RedisConfig
	DB    DBConfig
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// DBConfig 数据库配置
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "album_db"),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量作为整数，如果不存在或无法转换则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
