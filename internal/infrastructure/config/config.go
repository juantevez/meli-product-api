package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Type          string
	ProductsFile  string
	SellersFile   string
	ReviewsFile   string
	QuestionsFile string
}

type LoggerConfig struct {
	Level  string
	Format string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Type:          getEnv("DB_TYPE", "json"),
			ProductsFile:  getEnv("PRODUCTS_FILE", "./data/products.json"),
			SellersFile:   getEnv("SELLERS_FILE", "./data/sellers.json"),
			ReviewsFile:   getEnv("REVIEWS_FILE", "./data/reviews.json"),
			QuestionsFile: getEnv("QUESTIONS_FILE", "./data/questions.json"),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func (c *Config) Validate() error {
	// Add validation logic here
	if c.Server.Port == "" {
		log.Fatal("SERVER_PORT is required")
	}
	return nil
}
