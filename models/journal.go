package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JournalEntry struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	Tags      []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	Mood      string             `json:"mood,omitempty" bson:"mood,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Reflection struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EntryID     primitive.ObjectID `json:"entry_id" bson:"entry_id"`
	UserID      string             `json:"user_id" bson:"user_id"`
	Content     string             `json:"content" bson:"content"`
	Type        string             `json:"type" bson:"type"` // "insight", "summary", "analysis"
	Keywords    []string           `json:"keywords,omitempty" bson:"keywords,omitempty"`
	Sentiment   string             `json:"sentiment,omitempty" bson:"sentiment,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type CreateJournalRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags,omitempty"`
	Mood    string   `json:"mood,omitempty"`
}

type ReflectionRequest struct {
	EntryID string `json:"entry_id"`
	Type    string `json:"type,omitempty"` // defaults to "insight"
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
} 