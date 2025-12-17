package store

import (
	"apigo1/firebase"
	"apigo1/models"
	"context"
	"strconv"
	"time"

	"google.golang.org/api/iterator"
)

// FirestoreStore manages todos in Firestore
type FirestoreStore struct {
	collection string
	ctx        context.Context
}

// NewFirestoreStore creates a new FirestoreStore
func NewFirestoreStore(ctx context.Context) *FirestoreStore {
	return &FirestoreStore{
		collection: "todos",
		ctx:        ctx,
	}
}

// GetAll returns all todos
func (s *FirestoreStore) GetAll() []*models.Todo {
	var todos []*models.Todo

	iter := firebase.FirestoreClient.Collection(s.collection).Documents(s.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return todos
		}

		todo := &models.Todo{}
		if err := doc.DataTo(todo); err != nil {
			continue
		}
		// Parse ID from document ID
		if id, err := strconv.Atoi(doc.Ref.ID); err == nil {
			todo.ID = id
		}
		todos = append(todos, todo)
	}

	return todos
}

// GetByID returns a todo by ID
func (s *FirestoreStore) GetByID(id int) (*models.Todo, bool) {
	docRef := firebase.FirestoreClient.Collection(s.collection).Doc(strconv.Itoa(id))
	doc, err := docRef.Get(s.ctx)
	if err != nil {
		return nil, false
	}

	todo := &models.Todo{}
	if err := doc.DataTo(todo); err != nil {
		return nil, false
	}
	todo.ID = id
	return todo, true
}

// Create creates a new todo
func (s *FirestoreStore) Create(todo *models.Todo) *models.Todo {
	// Get the next ID by counting documents
	docs, err := firebase.FirestoreClient.Collection(s.collection).Documents(s.ctx).GetAll()
	if err != nil {
		// Fallback: use timestamp as ID
		todo.ID = int(time.Now().Unix())
	} else {
		maxID := 0
		for _, doc := range docs {
			if id, err := strconv.Atoi(doc.Ref.ID); err == nil && id > maxID {
				maxID = id
			}
		}
		todo.ID = maxID + 1
	}

	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	docRef := firebase.FirestoreClient.Collection(s.collection).Doc(strconv.Itoa(todo.ID))
	_, err = docRef.Set(s.ctx, todo)
	if err != nil {
		return nil
	}

	return todo
}

// Update updates an existing todo
func (s *FirestoreStore) Update(id int, updatedTodo *models.Todo) (*models.Todo, bool) {
	docRef := firebase.FirestoreClient.Collection(s.collection).Doc(strconv.Itoa(id))
	doc, err := docRef.Get(s.ctx)
	if err != nil {
		return nil, false
	}

	existingTodo := &models.Todo{}
	if err := doc.DataTo(existingTodo); err != nil {
		return nil, false
	}

	// Merge updates
	if updatedTodo.Title != "" {
		existingTodo.Title = updatedTodo.Title
	}
	if updatedTodo.Description != "" {
		existingTodo.Description = updatedTodo.Description
	}
	existingTodo.Completed = updatedTodo.Completed
	existingTodo.UpdatedAt = time.Now()

	_, err = docRef.Set(s.ctx, existingTodo)
	if err != nil {
		return nil, false
	}

	existingTodo.ID = id
	return existingTodo, true
}

// Delete deletes a todo by ID
func (s *FirestoreStore) Delete(id int) bool {
	docRef := firebase.FirestoreClient.Collection(s.collection).Doc(strconv.Itoa(id))
	_, err := docRef.Get(s.ctx)
	if err != nil {
		return false
	}

	_, err = docRef.Delete(s.ctx)
	return err == nil
}

