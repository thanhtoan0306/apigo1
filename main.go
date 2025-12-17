package main

import (
	"apigo1/firebase"
	"apigo1/handlers"
	"apigo1/store"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	ctx := context.Background()

	// Initialize Firebase
	if err := firebase.InitializeFirebase(ctx); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	defer firebase.Close()

	// Initialize Firestore store
	todoStore := store.NewFirestoreStore(ctx)

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoStore)

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	
	// Todo routes
	api.HandleFunc("/todos", todoHandler.GetAllTodos).Methods("GET")
	api.HandleFunc("/todos/{id}", todoHandler.GetTodoByID).Methods("GET")
	api.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	api.HandleFunc("/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
	api.HandleFunc("/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Start server
	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("API endpoints available at http://localhost%s/api", port)

	// Graceful shutdown
	go func() {
		if err := http.ListenAndServe(port, router); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

