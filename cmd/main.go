package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"soulprint-backend/config"
	"soulprint-backend/controllers"
	"soulprint-backend/routes"
	"soulprint-backend/services"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to MongoDB
	mongoClient, err := connectMongoDB()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Initialize services
	journalService := services.NewJournalService(mongoClient)
	aiService := services.NewAIService(mongoClient, journalService)

	// Initialize controllers
	journalController := controllers.NewJournalController(journalService)
	reflectionController := controllers.NewReflectionController(aiService)

	// Setup routes
	router := routes.NewRouter(journalController, reflectionController)

	// Start server
	port := config.AppConfig.Port
	fmt.Printf("🌟 Soulprint Backend starting on port %s\n", port)
	fmt.Printf("📖 MongoDB: %s\n", config.AppConfig.MongoDatabase)
	if config.AppConfig.UseLocalModel {
		fmt.Printf("🤖 AI Model: Local Llama3 (%s)\n", config.AppConfig.LocalModelURL)
	} else {
		fmt.Printf("🤖 AI Model: %s\n", config.AppConfig.OpenAIModel)
	}
	fmt.Println("✨ Available endpoints:")
	fmt.Println("   GET  /health")
	fmt.Println("   POST /api/v1/user")
	fmt.Println("   POST /api/v1/entries")
	fmt.Println("   GET  /api/v1/entries")
	fmt.Println("   GET  /api/v1/entries/{id}")
	fmt.Println("   PUT  /api/v1/entries/{id}")
	fmt.Println("   DELETE /api/v1/entries/{id}")
	fmt.Println("   POST /api/v1/reflect")
	fmt.Println("   GET  /api/v1/insights")
	fmt.Println("   GET  /api/v1/reflections")
	fmt.Println("   GET  /api/v1/entries/{id}/reflections")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func connectMongoDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoURI))
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Connected to MongoDB!")
	return client, nil
}
