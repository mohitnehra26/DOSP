package memory

import (
	"testing"

	"reddit-clone/internal/models"
)

func TestCreateAndGetUser(t *testing.T) {
	store := NewMemoryStore()
	user := &models.User{
		ID:       "user1",
		Username: "testuser",
		Password: "password123",
	}

	err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrievedUser, err := store.GetUser("user1")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if retrievedUser.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", retrievedUser.Username)
	}
}

func TestCreateAndGetPost(t *testing.T) {
	store := NewMemoryStore()
	post := &models.Post{
		ID:          "post1",
		SubredditID: "subreddit1",
		AuthorID:    "user1",
		Title:       "Test Post",
		Content:     "This is a test post",
	}

	err := store.CreatePost(post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	retrievedPost, err := store.GetPost("post1")
	if err != nil {
		t.Fatalf("Failed to get post: %v", err)
	}

	if retrievedPost.Title != "Test Post" {
		t.Errorf("Expected post title 'Test Post', got '%s'", retrievedPost.Title)
	}
}
