package models

import "time"

// Blog represents a blog post
type Blog struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`     // Markdown content
	Slug        string    `json:"slug"`        // URL-friendly identifier
	Author      string    `json:"author"`
	Published   bool      `json:"published"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateBlogRequest represents the request body for creating a blog
type CreateBlogRequest struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Slug      string   `json:"slug"`
	Author    string   `json:"author"`
	Published bool     `json:"published"`
	Tags      []string `json:"tags"`
}

// UpdateBlogRequest represents the request body for updating a blog
type UpdateBlogRequest struct {
	Title     *string   `json:"title"`
	Content   *string   `json:"content"`
	Slug      *string   `json:"slug"`
	Author    *string   `json:"author"`
	Published *bool     `json:"published"`
	Tags      *[]string `json:"tags"`
}

