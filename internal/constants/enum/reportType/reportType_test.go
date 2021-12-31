package reportType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValid(t *testing.T) {
	assert.True(t, IsValidReport("members"))
	assert.False(t, IsValidReport("xx"))
}
