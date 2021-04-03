package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryFilter(t *testing.T) {
	t.Run("Contains Name filter", func(t *testing.T) {
		filter := QueryFilters{}
		filter.AddFilter("Name", "Bruno")
		assert.Equal(t, "Bruno", filter.GetFilter("Name").(string))
	})
	t.Run("Contains no Name filter", func(t *testing.T) {
		filter := QueryFilters{}
		assert.Nil(t, filter.GetFilter("Name"))
	})
}