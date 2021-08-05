package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestID(t *testing.T) {
	assert.True(t, IsValidID(NewID()))
	assert.False(t, IsValidID(""))
}
