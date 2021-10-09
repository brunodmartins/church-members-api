package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordValidation(t *testing.T) {
	assert.False(t, isValidPassword("12345"))
	assert.False(t, isValidPassword("123456"))
	assert.False(t, isValidPassword("12345a"))
	assert.False(t, isValidPassword("1234aB"))
	assert.True(t, isValidPassword("1234aB@"))
}
