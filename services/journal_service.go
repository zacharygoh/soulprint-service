package services

import (
	"context"
	"fmt"
	"time"

	"soulprint-backend/config"
	"soulprint-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JournalService struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewJournalService(client *mongo.Client) *JournalService {
	collection := client.Database(config.AppConfig.MongoDatabase).Collection("journal_entries")
	return &JournalService{
		client:     client,
		collection: collection,
	}
}

func (js *JournalService) CreateEntry(userID string, req models.CreateJournalRequest) (*models.JournalEntry, error) {
	entry := &models.JournalEntry{
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		Mood:      req.Mood,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := js.collection.InsertOne(context.Background(), entry)
	if err != nil {
		return nil, fmt.Errorf("failed to create journal entry: %w", err)
	}

	entry.ID = result.InsertedID.(primitive.ObjectID)
	return entry, nil
}

func (js *JournalService) GetEntries(userID string) ([]models.JournalEntry, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := js.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find journal entries: %w", err)
	}
	defer cursor.Close(context.Background())

	var entries []models.JournalEntry
	if err = cursor.All(context.Background(), &entries); err != nil {
		return nil, fmt.Errorf("failed to decode journal entries: %w", err)
	}

	return entries, nil
}

func (js *JournalService) GetEntryByID(userID, entryID string) (*models.JournalEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		return nil, fmt.Errorf("invalid entry ID: %w", err)
	}

	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	var entry models.JournalEntry
	err = js.collection.FindOne(context.Background(), filter).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, fmt.Errorf("failed to find journal entry: %w", err)
	}

	return &entry, nil
}

func (js *JournalService) UpdateEntry(userID, entryID string, req models.CreateJournalRequest) (*models.JournalEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		return nil, fmt.Errorf("invalid entry ID: %w", err)
	}

	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	update := bson.M{
		"$set": bson.M{
			"title":      req.Title,
			"content":    req.Content,
			"tags":       req.Tags,
			"mood":       req.Mood,
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var entry models.JournalEntry
	err = js.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, fmt.Errorf("failed to update journal entry: %w", err)
	}

	return &entry, nil
}

func (js *JournalService) DeleteEntry(userID, entryID string) error {
	objectID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		return fmt.Errorf("invalid entry ID: %w", err)
	}

	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	result, err := js.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete journal entry: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("journal entry not found")
	}

	return nil
} 