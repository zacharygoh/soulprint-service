package services

import (
	"context"
	"fmt"
	"time"

	"soulprint-backend/config"
	"soulprint-backend/models"
	"soulprint-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AIService struct {
	client           *mongo.Client
	collection       *mongo.Collection
	journalService   *JournalService
	openaiClient     *utils.OpenAIClient
}

func NewAIService(client *mongo.Client, journalService *JournalService) *AIService {
	collection := client.Database(config.AppConfig.MongoDatabase).Collection("reflections")
	return &AIService{
		client:         client,
		collection:     collection,
		journalService: journalService,
		openaiClient:   utils.NewOpenAIClient(),
	}
}

func (ais *AIService) GenerateReflection(userID string, req models.ReflectionRequest) (*models.Reflection, error) {
	// Get the journal entry
	entry, err := ais.journalService.GetEntryByID(userID, req.EntryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal entry: %w", err)
	}

	// Default reflection type
	reflectionType := req.Type
	if reflectionType == "" {
		reflectionType = "insight"
	}

	// Generate AI reflection
	reflectionContent, err := ais.openaiClient.GenerateReflection(entry.Content, reflectionType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI reflection: %w", err)
	}

	// Extract keywords (optional, can fail gracefully)
	keywords, _ := ais.openaiClient.ExtractKeywords(entry.Content)

	// Create reflection record
	reflection := &models.Reflection{
		EntryID:   entry.ID,
		UserID:    userID,
		Content:   reflectionContent,
		Type:      reflectionType,
		Keywords:  keywords,
		Sentiment: ais.extractSentiment(reflectionContent), // Simple sentiment analysis
		CreatedAt: time.Now(),
	}

	result, err := ais.collection.InsertOne(context.Background(), reflection)
	if err != nil {
		return nil, fmt.Errorf("failed to save reflection: %w", err)
	}

	reflection.ID = result.InsertedID.(primitive.ObjectID)
	return reflection, nil
}

func (ais *AIService) GetReflections(userID string) ([]models.Reflection, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := ais.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find reflections: %w", err)
	}
	defer cursor.Close(context.Background())

	var reflections []models.Reflection
	if err = cursor.All(context.Background(), &reflections); err != nil {
		return nil, fmt.Errorf("failed to decode reflections: %w", err)
	}

	return reflections, nil
}

func (ais *AIService) GetReflectionsByEntry(userID, entryID string) ([]models.Reflection, error) {
	objectID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		return nil, fmt.Errorf("invalid entry ID: %w", err)
	}

	filter := bson.M{
		"entry_id": objectID,
		"user_id":  userID,
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := ais.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find reflections: %w", err)
	}
	defer cursor.Close(context.Background())

	var reflections []models.Reflection
	if err = cursor.All(context.Background(), &reflections); err != nil {
		return nil, fmt.Errorf("failed to decode reflections: %w", err)
	}

	return reflections, nil
}

func (ais *AIService) GetInsights(userID string) (map[string]interface{}, error) {
	// Get recent reflections
	reflections, err := ais.GetReflections(userID)
	if err != nil {
		return nil, err
	}

	// Simple insights aggregation
	insights := map[string]interface{}{
		"total_reflections": len(reflections),
		"recent_themes":     ais.extractRecentThemes(reflections),
		"sentiment_trends":  ais.analyzeSentimentTrends(reflections),
		"reflection_types":  ais.countReflectionTypes(reflections),
	}

	return insights, nil
}

// Helper methods
func (ais *AIService) extractSentiment(content string) string {
	// Simple sentiment analysis based on keywords
	// In production, you might want to use a proper sentiment analysis library or API
	positive := []string{"happy", "joy", "grateful", "excited", "love", "wonderful", "amazing", "great"}
	negative := []string{"sad", "angry", "frustrated", "worried", "anxious", "terrible", "awful", "horrible"}
	
	positiveCount := 0
	negativeCount := 0
	
	contentLower := content
	for _, word := range positive {
		if contains(contentLower, word) {
			positiveCount++
		}
	}
	
	for _, word := range negative {
		if contains(contentLower, word) {
			negativeCount++
		}
	}
	
	if positiveCount > negativeCount {
		return "positive"
	} else if negativeCount > positiveCount {
		return "negative"
	}
	return "neutral"
}

func (ais *AIService) extractRecentThemes(reflections []models.Reflection) []string {
	themeMap := make(map[string]int)
	
	// Collect keywords from recent reflections (last 10)
	limit := 10
	if len(reflections) < limit {
		limit = len(reflections)
	}
	
	for i := 0; i < limit; i++ {
		for _, keyword := range reflections[i].Keywords {
			themeMap[keyword]++
		}
	}
	
	// Return top themes
	var themes []string
	for theme := range themeMap {
		themes = append(themes, theme)
		if len(themes) >= 5 {
			break
		}
	}
	
	return themes
}

func (ais *AIService) analyzeSentimentTrends(reflections []models.Reflection) map[string]int {
	sentiments := map[string]int{
		"positive": 0,
		"negative": 0,
		"neutral":  0,
	}
	
	for _, reflection := range reflections {
		if sentiment := reflection.Sentiment; sentiment != "" {
			sentiments[sentiment]++
		}
	}
	
	return sentiments
}

func (ais *AIService) countReflectionTypes(reflections []models.Reflection) map[string]int {
	types := make(map[string]int)
	
	for _, reflection := range reflections {
		types[reflection.Type]++
	}
	
	return types
}

// Helper function to check if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
} 