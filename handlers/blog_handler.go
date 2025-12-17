package handlers

import (
	"apigo1/models"
	"apigo1/store"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// BlogHandler handles blog-related HTTP requests
type BlogHandler struct {
	store *store.BlogStore
}

// NewBlogHandler creates a new BlogHandler
func NewBlogHandler(s *store.BlogStore) *BlogHandler {
	return &BlogHandler{store: s}
}

// GetAllBlogs handles GET /blogs
func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	blogs := h.store.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    blogs,
	})
}

// GetBlogByID handles GET /blogs/{id}
func (h *BlogHandler) GetBlogByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	blog, exists := h.store.GetByID(id)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Blog not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    blog,
	})
}

// GetBlogBySlug handles GET /blogs/slug/{slug}
func (h *BlogHandler) GetBlogBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	blog, exists := h.store.GetBySlug(slug)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Blog not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    blog,
	})
}

// CreateBlog handles POST /blogs
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBlogRequest
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

	// Generate slug from title if not provided
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Title)
	}

	now := time.Now()
	blog := &models.Blog{
		Title:     req.Title,
		Content:   req.Content,
		Slug:      slug,
		Author:    req.Author,
		Published: req.Published,
		Tags:      req.Tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	createdBlog := h.store.Create(blog)
	if createdBlog == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to create blog",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    createdBlog,
	})
}

// UpdateBlog handles PUT /blogs/{id}
func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingBlog, exists := h.store.GetByID(id)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Blog not found",
		})
		return
	}

	updatedBlog := &models.Blog{
		ID:        existingBlog.ID,
		Title:     existingBlog.Title,
		Content:   existingBlog.Content,
		Slug:      existingBlog.Slug,
		Author:    existingBlog.Author,
		Published: existingBlog.Published,
		Tags:      existingBlog.Tags,
		UpdatedAt: time.Now(),
	}

	// Update fields if provided
	if req.Title != nil {
		updatedBlog.Title = *req.Title
	}
	if req.Content != nil {
		updatedBlog.Content = *req.Content
	}
	if req.Slug != nil {
		updatedBlog.Slug = *req.Slug
	} else if req.Title != nil {
		// Regenerate slug if title changed but slug not provided
		updatedBlog.Slug = generateSlug(*req.Title)
	}
	if req.Author != nil {
		updatedBlog.Author = *req.Author
	}
	if req.Published != nil {
		updatedBlog.Published = *req.Published
	}
	if req.Tags != nil {
		updatedBlog.Tags = *req.Tags
	}

	blog, _ := h.store.Update(id, updatedBlog)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    blog,
	})
}

// DeleteBlog handles DELETE /blogs/{id}
func (h *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	exists := h.store.Delete(id)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Blog not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message":  "Blog deleted successfully",
	})
}

// generateSlug generates a URL-friendly slug from title
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.TrimSpace(slug)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}
	slug = result.String()
	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	return slug
}

