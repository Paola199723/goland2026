package entities

import (
	"testing"
)

func TestUserCreation(t *testing.T) {
	user := User{Name: "John Doe", Email: "john@example.com"}
	if user.Name != "John Doe" {
		t.Errorf("Expected Name to be 'John Doe', got '%s'", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected Email to be 'john@example.com', got '%s'", user.Email)
	}
}