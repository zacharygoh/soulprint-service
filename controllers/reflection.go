package controllers

import (
	"encoding/json"
	"net/http"

	"soulprint-backend/models"
	"soulprint-backend/services"

	"github.com/gorilla/mux"
)

type ReflectionController struct {
	aiService *services.AIService
}

func NewReflectionController(aiService *services.AIService) *ReflectionController {
	return &ReflectionController{
		aiService: aiService,
	}
}

// POST /reflect
func (rc *ReflectionController) GenerateReflection(w http.ResponseWriter, r *http.Request) {
	var req models.ReflectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.EntryID == "" {
		http.Error(w, "Entry ID is required", http.StatusBadRequest)
		return
	}

	// For MVP, use hardcoded user ID
	userID := "user123"
	
	reflection, err := rc.aiService.GenerateReflection(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    reflection,
	})
}

// GET /insights
func (rc *ReflectionController) GetInsights(w http.ResponseWriter, r *http.Request) {
	// For MVP, use hardcoded user ID
	userID := "user123"
	
	insights, err := rc.aiService.GetInsights(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    insights,
	})
}

// GET /reflections
func (rc *ReflectionController) GetReflections(w http.ResponseWriter, r *http.Request) {
	// For MVP, use hardcoded user ID
	userID := "user123"
	
	reflections, err := rc.aiService.GetReflections(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    reflections,
	})
}

// GET /entries/{id}/reflections
func (rc *ReflectionController) GetReflectionsByEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entryID := vars["id"]
	
	if entryID == "" {
		http.Error(w, "Entry ID is required", http.StatusBadRequest)
		return
	}

	// For MVP, use hardcoded user ID
	userID := "user123"
	
	reflections, err := rc.aiService.GetReflectionsByEntry(userID, entryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    reflections,
	})
}

// POST /user (MVP hardcoded user creation)
func (rc *ReflectionController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For MVP, return hardcoded user response
	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    user,
		"message": "User created successfully (MVP hardcoded)",
	})
} 