package handlers

import (
	"apigo1/models"
	"apigo1/store"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// TodoHandler handles todo-related HTTP requests
type TodoHandler struct {
	store store.TodoStoreInterface
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(s store.TodoStoreInterface) *TodoHandler {
	return &TodoHandler{store: s}
}

// GetAllTodos handles GET /todos
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos := h.store.GetAll()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    todos,
	})
}

// GetTodoByID handles GET /todos/{id}
func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, exists := h.store.GetByID(id)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Todo not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    todo,
	})
}

// CreateTodo handles POST /todos
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Title is required",
		})
		return
	}

	now := time.Now()
	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	createdTodo := h.store.Create(todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    createdTodo,
	})
}

// UpdateTodo handles PUT /todos/{id}
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingTodo, exists := h.store.GetByID(id)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Todo not found",
		})
		return
	}

	updatedTodo := &models.Todo{
		ID:          existingTodo.ID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   existingTodo.Completed,
		UpdatedAt:   time.Now(),
	}

	if req.Completed != nil {
		updatedTodo.Completed = *req.Completed
	}

	if req.Title == "" {
		updatedTodo.Title = existingTodo.Title
	}
	if req.Description == "" {
		updatedTodo.Description = existingTodo.Description
	}

	todo, _ := h.store.Update(id, updatedTodo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    todo,
	})
}

// DeleteTodo handles DELETE /todos/{id}
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	exists := h.store.Delete(id)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Todo not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message":  "Todo deleted successfully",
	})
}

