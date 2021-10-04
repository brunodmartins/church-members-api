package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertDate(t *testing.T) {
	assert.Equal(t, "01-01", ConvertDate(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.Equal(t, "12-24", ConvertDate(time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC)))
}

