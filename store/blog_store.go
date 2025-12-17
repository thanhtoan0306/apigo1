package store

import (
	"apigo1/firebase"
	"apigo1/models"
	"context"
	"strconv"
	"time"

	"google.golang.org/api/iterator"
)

// BlogStore manages blogs in Firestore
type BlogStore struct {
	collection string
	ctx        context.Context
}

// NewBlogStore creates a new BlogStore
func NewBlogStore(ctx context.Context) *BlogStore {
	return &BlogStore{
		collection: "blogs",
		ctx:        ctx,
	}
}

// GetAll returns all blogs
func (s *BlogStore) GetAll() []*models.Blog {
	var blogs []*models.Blog

	iter := firebase.FirestoreClient.Collection(s.collection).Documents(s.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return blogs
		}

		blog := &models.Blog{}
		if err := doc.DataTo(blog); err != nil {
			continue
		}
		// Parse ID from document ID
		if id, err := strconv.Atoi(doc.Ref.ID); err == nil {
			blog.ID = id
		}
		blogs = append(blogs, blog)
	}

	return blogs
}

// GetByID returns a blog by ID
func (s *BlogStore) GetByID(id int) (*models.Blog, bool) {
	docRef := firebase.FirestoreClient.Collection(s.collection).Doc(strconv.Itoa(id))
	doc, err := docRef.Get(s.ctx)
	if err != nil {
		return nil, false
	}

	blog := &models.Blog{}
	if err := doc.DataTo(blog); err != nil {
		return nil, false
	}
	blog.ID = id
	return blog, true
}

// GetBySlug returns a blog by slug
func (s *BlogStore) GetBySlug(slug string) (*models.Blog, bool) {
	iter := firebase.FirestoreClient.Collection(s.collection).Where("slug", "==", slug).Documents(s.ctx)
	docs, err := iter.GetAll()
	if err != nil || len(docs) == 0 {
		return nil, false
	}

	blog := &models.Blog{}
	if err := docs[0].DataTo(blog); err != nil {
		return nil, false
	}
	if id, err := strconv.Atoi(docs[0].Ref.ID); err == nil {
		blog.ID = id
	}
	return blog, true
}

// Create creates a new blog
func (s *BlogStore) Create(blog *models.Blog) *models.Blog {
	// Get the next ID by counting documents
	docs, err := firebase.FirestoreClient.Collection(s.collection).Documents(s.ctx).GetAll()
	if err != nil {
		// Fallback: use timestamp as ID
		blog.ID = int(time.Now().Unix())
	} else {
		maxID := 0
		for _, doc := range docs {
			if id, err := strconv.Atoi(doc.Ref.ID); err == nil && id > maxID {
				maxID = id
			}
		}
		blog.ID = maxID + 1
	}

	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now

	docRef := firebase.FirestoreClient.Collection(s.collection).Doc(strconv.Itoa(blog.ID))
	_, err = docRef.Set(s.ctx, blog)
	if err != nil {
		return nil
	}

	return blog
}

// Update updates an existing blog
func (s *BlogStore) Update(id int, updatedBlog *models.Blog) (*models.Blog, bool) {
	docRef := firebase.FirestoreClient.Collection(s.collection).Doc(strconv.Itoa(id))
	doc, err := docRef.Get(s.ctx)
	if err != nil {
		return nil, false
	}

	existingBlog := &models.Blog{}
	if err := doc.DataTo(existingBlog); err != nil {
		return nil, false
	}

	// Merge updates
	if updatedBlog.Title != "" {
		existingBlog.Title = updatedBlog.Title
	}
	if updatedBlog.Content != "" {
		existingBlog.Content = updatedBlog.Content
	}
	if updatedBlog.Slug != "" {
		existingBlog.Slug = updatedBlog.Slug
	}
	if updatedBlog.Author != "" {
		existingBlog.Author = updatedBlog.Author
	}
	existingBlog.Published = updatedBlog.Published
	if updatedBlog.Tags != nil {
		existingBlog.Tags = updatedBlog.Tags
	}
	existingBlog.UpdatedAt = time.Now()

	_, err = docRef.Set(s.ctx, existingBlog)
	if err != nil {
		return nil, false
	}

	existingBlog.ID = id
	return existingBlog, true
}

// Delete deletes a blog by ID
func (s *BlogStore) Delete(id int) bool {
	docRef := firebase.FirestoreClient.Collection(s.collection).Doc(strconv.Itoa(id))
	_, err := docRef.Get(s.ctx)
	if err != nil {
		return false
	}

	_, err = docRef.Delete(s.ctx)
	return err == nil
}

