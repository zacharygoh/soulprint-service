package routes

import (
	"net/http"

	"soulprint-backend/controllers"

	"github.com/gorilla/mux"
)

func NewRouter(journalController *controllers.JournalController, reflectionController *controllers.ReflectionController) *mux.Router {
	router := mux.NewRouter()

	// Add CORS middleware
	router.Use(corsMiddleware)

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// User routes (MVP)
	api.HandleFunc("/user", reflectionController.CreateUser).Methods("POST")

	// Journal entry routes
	api.HandleFunc("/entries", journalController.CreateEntry).Methods("POST")
	api.HandleFunc("/entries", journalController.GetEntries).Methods("GET")
	api.HandleFunc("/entries/{id}", journalController.GetEntry).Methods("GET")
	api.HandleFunc("/entries/{id}", journalController.UpdateEntry).Methods("PUT")
	api.HandleFunc("/entries/{id}", journalController.DeleteEntry).Methods("DELETE")

	// AI reflection routes
	api.HandleFunc("/reflect", reflectionController.GenerateReflection).Methods("POST")
	api.HandleFunc("/insights", reflectionController.GetInsights).Methods("GET")
	api.HandleFunc("/reflections", reflectionController.GetReflections).Methods("GET")
	api.HandleFunc("/entries/{id}/reflections", reflectionController.GetReflectionsByEntry).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "soulprint-backend"}`))
	}).Methods("GET")

	// Root endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Welcome to Soulprint API", "version": "1.0.0"}`))
	}).Methods("GET")

	return router
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
