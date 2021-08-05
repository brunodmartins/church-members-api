package member

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuerySpecification(t *testing.T) {
	t.Run("Contains Name filter", func(t *testing.T) {
		filter := new(QuerySpecification)
		filter.AddFilter("Name", "Bruno")
		assert.Equal(t, "Bruno", filter.values["Name"].(string))
	})
	t.Run("Contains no Name filter", func(t *testing.T) {
		filter := new(QuerySpecification)
		assert.Nil(t, filter.values["Name"])
	})
}


func TestCreateActiveFilter(t *testing.T) {
	builder := expression.NewBuilder()
	spec := CreateActiveFilter()
	builder = spec(builder)
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 1)
}

func TestQuerySpecification_ApplyFilters(t *testing.T) {
	assertFilters := func(querySpec *QuerySpecification, length int) {
		builder := expression.NewBuilder()
		spec := querySpec.ToSpecification()
		builder = spec(builder)
		expression, _ := builder.Build()
		assert.Len(t, expression.Names(), length)
	}
	t.Run("Without filters", func(t *testing.T) {
		spec := new(QuerySpecification)
		assertFilters(spec, 0)
	})
	t.Run("With one filter", func(t *testing.T) {
		spec := new(QuerySpecification)
		spec.AddFilter("name", "test")
		assertFilters(spec, 1)
	})
	t.Run("With two filter", func(t *testing.T) {
		spec := new(QuerySpecification)
		spec.AddFilter("name", "test")
		spec.AddFilter("active", true)
		assertFilters(spec, 2)
	})
	t.Run("With three filter", func(t *testing.T) {
		spec := new(QuerySpecification)
		spec.AddFilter("name", "test")
		spec.AddFilter("active", true)
		spec.AddFilter("gender", "M")
		assertFilters(spec, 3)
	})
}