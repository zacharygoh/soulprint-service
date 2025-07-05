package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	MongoURI       string
	MongoDatabase  string
	OpenAIAPIKey   string
	OpenAIModel    string
	LocalModelURL  string
	LocalModelName string
	UseLocalModel  bool
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	AppConfig = &Config{
		Port:           getEnv("PORT", "8080"),
		MongoURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase:  getEnv("MONGODB_DATABASE", "soulprint"),
		OpenAIAPIKey:   getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:    getEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
		LocalModelURL:  getEnv("LOCAL_MODEL_URL", "http://localhost:11434"),
		LocalModelName: getEnv("LOCAL_MODEL_NAME", "llama3"),
		UseLocalModel:  getEnv("USE_LOCAL_MODEL", "false") == "true",
	}

	if AppConfig.UseLocalModel {
		// Local model configuration loaded
	} else if AppConfig.OpenAIAPIKey == "" {
		log.Println("Warning: OPENAI_API_KEY not set and USE_LOCAL_MODEL is false")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
