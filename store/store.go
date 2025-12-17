package store

import (
	"apigo1/models"
	"sync"
)

// TodoStore manages todos in memory
type TodoStore struct {
	todos map[int]*models.Todo
	mu    sync.RWMutex
	nextID int
}

// NewTodoStore creates a new TodoStore
func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos:  make(map[int]*models.Todo),
		nextID: 1,
	}
}

// GetAll returns all todos
func (s *TodoStore) GetAll() []*models.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]*models.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos
}

// GetByID returns a todo by ID
func (s *TodoStore) GetByID(id int) (*models.Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, exists := s.todos[id]
	return todo, exists
}

// Create creates a new todo
func (s *TodoStore) Create(todo *models.Todo) *models.Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo.ID = s.nextID
	s.nextID++
	s.todos[todo.ID] = todo
	return todo
}

// Update updates an existing todo
func (s *TodoStore) Update(id int, updatedTodo *models.Todo) (*models.Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, exists := s.todos[id]
	if !exists {
		return nil, false
	}

	if updatedTodo.Title != "" {
		todo.Title = updatedTodo.Title
	}
	if updatedTodo.Description != "" {
		todo.Description = updatedTodo.Description
	}
	todo.Completed = updatedTodo.Completed
	todo.UpdatedAt = updatedTodo.UpdatedAt

	return todo, true
}

// Delete deletes a todo by ID
func (s *TodoStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.todos[id]
	if exists {
		delete(s.todos, id)
	}
	return exists
}

