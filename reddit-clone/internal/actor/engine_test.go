package actor

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	pb "reddit-clone/api/proto/generated"
	"reddit-clone/internal/models"
	"reddit-clone/internal/store/memory"
	"reddit-clone/pkg/metrics"
)

func TestHandleGetDirectMessages(t *testing.T) {
	// Initialize the memory store
	store := memory.NewMemoryStore()

	// Add test data
	store.SendMessage(&models.DirectMessage{
		ID:        "msg1",
		FromID:    "user1",
		ToID:      "user2",
		Content:   "Hello!",
		Timestamp: time.Now().Unix(),
	})

	// Create the engine actor
	engine := NewEngineActor(store, metrics.NewRedditMetrics())
	props := actor.PropsFromProducer(func() actor.Actor { return engine })

	// Create a root context and spawn the engine actor
	system := actor.NewActorSystem()
	rootContext := system.Root
	enginePID, err := rootContext.SpawnNamed(props, "engine")
	if err != nil {
		t.Fatalf("Failed to spawn engine actor: %v", err)
	}

	// Send GetDirectMessagesMessage to the engine actor
	future := rootContext.RequestFuture(enginePID, &pb.GetDirectMessagesMessage{UserId: "user2"}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		t.Fatalf("Failed to get response from engine actor: %v", err)
	}

	// Assert the response type and content
	response, ok := result.(*pb.DirectMessagesResponse)
	if !ok {
		t.Fatalf("Expected DirectMessagesResponse, got %T", result)
	}

	if len(response.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(response.Messages))
	}

	if response.Messages[0].Content != "Hello!" {
		t.Errorf("Expected message content 'Hello!', got '%s'", response.Messages[0].Content)
	}
}

func TestHandleUserMessage(t *testing.T) {
	system := actor.NewActorSystem()
	store := memory.NewMemoryStore()
	metrics := metrics.NewRedditMetrics()
	engine := NewEngineActor(store, metrics)

	props := actor.PropsFromProducer(func() actor.Actor { return engine })
	enginePID, err := system.Root.SpawnNamed(props, "engine")
	if err != nil {
		t.Fatalf("Failed to spawn engine actor: %v", err)
	}

	userMsg := &pb.UserMessage{
		UserId:   "user1",
		Username: "testuser",
		Password: "password123",
	}

	future := system.Root.RequestFuture(enginePID, userMsg, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		t.Fatalf("Failed to get response from engine actor: %v", err)
	}

	response, ok := result.(*pb.SuccessResponse)
	if !ok {
		t.Fatalf("Expected SuccessResponse, got %T", result)
	}
	if response.Message != "User registered successfully" {
		t.Errorf("Expected 'User registered successfully', got '%s'", response.Message)
	}

	user, err := store.GetUser("user1")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
}

func TestHandlePostMessage(t *testing.T) {
	system := actor.NewActorSystem()
	store := memory.NewMemoryStore()
	metrics := metrics.NewRedditMetrics()
	engine := NewEngineActor(store, metrics)

	props := actor.PropsFromProducer(func() actor.Actor { return engine })
	enginePID, err := system.Root.SpawnNamed(props, "engine")
	if err != nil {
		t.Fatalf("Failed to spawn engine actor: %v", err)
	}

	postMsg := &pb.PostMessage{
		Id:          "post1",
		SubredditId: "subreddit1",
		AuthorId:    "user1",
		Title:       "Test Post",
		Content:     "This is a test post",
	}

	future := system.Root.RequestFuture(enginePID, postMsg, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		t.Fatalf("Failed to get response from engine actor: %v", err)
	}

	response, ok := result.(*pb.SuccessResponse)
	if !ok {
		t.Fatalf("Expected SuccessResponse, got %T", result)
	}
	if response.Message != "Post created successfully" {
		t.Errorf("Expected 'Post created successfully', got '%s'", response.Message)
	}

	post, err := store.GetPost("post1")
	if err != nil {
		t.Fatalf("Failed to get post: %v", err)
	}
	if post.Title != "Test Post" {
		t.Errorf("Expected post title 'Test Post', got '%s'", post.Title)
	}
}
