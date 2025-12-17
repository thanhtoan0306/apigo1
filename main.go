package main

import (
	"apigo1/docs"
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
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Todo & Blog API
// @version         1.0
// @description     Backend API cho ứng dụng Todo List và Blog Management với Firebase Firestore. Hỗ trợ Markdown content cho blogs.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @BasePath  /api
//
// @schemes   http https
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

	// Initialize Firestore stores
	todoStore := store.NewFirestoreStore(ctx)
	blogStore := store.NewBlogStore(ctx)

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoStore)
	blogHandler := handlers.NewBlogHandler(blogStore)

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

	// Blog routes
	api.HandleFunc("/blogs", blogHandler.GetAllBlogs).Methods("GET")
	api.HandleFunc("/blogs/{id}", blogHandler.GetBlogByID).Methods("GET")
	api.HandleFunc("/blogs/slug/{slug}", blogHandler.GetBlogBySlug).Methods("GET")
	api.HandleFunc("/blogs", blogHandler.CreateBlog).Methods("POST")
	api.HandleFunc("/blogs/{id}", blogHandler.UpdateBlog).Methods("PUT")
	api.HandleFunc("/blogs/{id}", blogHandler.DeleteBlog).Methods("DELETE")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Swagger documentation
	// Mặc định Host rỗng để Swagger dùng origin hiện tại (production / local).
	docs.SwaggerInfo.Host = ""
	if envHost := os.Getenv("HOST"); envHost != "" {
		docs.SwaggerInfo.Host = envHost
	}
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

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
