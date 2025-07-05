package controllers

import (
	"encoding/json"
	"net/http"

	"soulprint-backend/models"
	"soulprint-backend/services"

	"github.com/gorilla/mux"
)

type JournalController struct {
	journalService *services.JournalService
}

func NewJournalController(journalService *services.JournalService) *JournalController {
	return &JournalController{
		journalService: journalService,
	}
}

// POST /entries
func (jc *JournalController) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var req models.CreateJournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" || req.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// For MVP, use hardcoded user ID
	userID := "user123"
	
	entry, err := jc.journalService.CreateEntry(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    entry,
	})
}

// GET /entries
func (jc *JournalController) GetEntries(w http.ResponseWriter, r *http.Request) {
	// For MVP, use hardcoded user ID
	userID := "user123"
	
	entries, err := jc.journalService.GetEntries(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    entries,
	})
}

// GET /entries/{id}
func (jc *JournalController) GetEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entryID := vars["id"]
	
	if entryID == "" {
		http.Error(w, "Entry ID is required", http.StatusBadRequest)
		return
	}

	// For MVP, use hardcoded user ID
	userID := "user123"
	
	entry, err := jc.journalService.GetEntryByID(userID, entryID)
	if err != nil {
		if err.Error() == "journal entry not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    entry,
	})
}

// PUT /entries/{id}
func (jc *JournalController) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entryID := vars["id"]
	
	if entryID == "" {
		http.Error(w, "Entry ID is required", http.StatusBadRequest)
		return
	}

	var req models.CreateJournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" || req.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// For MVP, use hardcoded user ID
	userID := "user123"
	
	entry, err := jc.journalService.UpdateEntry(userID, entryID, req)
	if err != nil {
		if err.Error() == "journal entry not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    entry,
	})
}

// DELETE /entries/{id}
func (jc *JournalController) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entryID := vars["id"]
	
	if entryID == "" {
		http.Error(w, "Entry ID is required", http.StatusBadRequest)
		return
	}

	// For MVP, use hardcoded user ID
	userID := "user123"
	
	err := jc.journalService.DeleteEntry(userID, entryID)
	if err != nil {
		if err.Error() == "journal entry not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Entry deleted successfully",
	})
} 