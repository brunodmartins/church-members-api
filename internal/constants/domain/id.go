package domain

import (
	"github.com/google/uuid"
)

// IsValidID validates the ID regex for UUID
func IsValidID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// NewID builds a new UUID
func NewID() string {
	return uuid.NewString()
}
