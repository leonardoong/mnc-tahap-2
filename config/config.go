package config

import (
	"log"
	"os"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	// Server
	ServerPort string

	// Database
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	CachePool     *redis.Pool
}

func LoadConfig() *Config {
	config := &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DBDriver:       getEnv("DB_DRIVER", "mysql"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "user"),
		DBPassword:     getEnv("DB_PASSWORD", "password"),
		DBName:         getEnv("DB_NAME", "emoney"),
		JWTSecret:      getEnv("JWT_SECRET", "jwt-secret"),
		JWTExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		RedisHost:      getEnv("REDIS_HOST", "redis"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Error converting %s to int, using default value %d: %v", key, defaultValue, err)
		return defaultValue
	}
	return value
}
