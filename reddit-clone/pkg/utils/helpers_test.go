package utils

import (
	"regexp"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id := GenerateID()
	if len(id) != 16 {
		t.Errorf("Generated ID length is incorrect. Expected 16, got %d", len(id))
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", id)
	if !matched {
		t.Errorf("Generated ID contains invalid characters: %s", id)
	}
}

func TestGenerateRandomContent(t *testing.T) {
	content := GenerateRandomContent()
	if content == "" {
		t.Error("Generated content is empty")
	}
}

func TestGenerateRandomTitle(t *testing.T) {
	title := GenerateRandomTitle()
	if title == "" {
		t.Error("Generated title is empty")
	}
}

func TestGenerateUsername(t *testing.T) {
	username := GenerateUsername(123)
	expected := "user_123"
	if username != expected {
		t.Errorf("Generated username is incorrect. Expected %s, got %s", expected, username)
	}
}
