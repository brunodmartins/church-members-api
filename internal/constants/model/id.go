package model

import (
	"github.com/google/uuid"
	"regexp"
)

//IsValidID validates the ID regex for UUID
func IsValidID(id string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(id)
}

//NewID builds a new UUID
func NewID() string {
	return uuid.NewString()
}
