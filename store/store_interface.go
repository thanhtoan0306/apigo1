package store

import "apigo1/models"

// TodoStoreInterface defines the interface for todo storage
type TodoStoreInterface interface {
	GetAll() []*models.Todo
	GetByID(id int) (*models.Todo, bool)
	Create(todo *models.Todo) *models.Todo
	Update(id int, updatedTodo *models.Todo) (*models.Todo, bool)
	Delete(id int) bool
}

