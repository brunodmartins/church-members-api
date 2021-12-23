package member

import (
	"context"
	"testing"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum/classification"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

	"github.com/stretchr/testify/assert"
)

func TestQuerySpecification(t *testing.T) {
	t.Run("Contains Name filter", func(t *testing.T) {
		filter := new(QueryBuilder)
		filter.AddFilter("Name", "Bruno")
		assert.Equal(t, "Bruno", filter.values["Name"].(string))
	})
	t.Run("Contains no Name filter", func(t *testing.T) {
		filter := new(QueryBuilder)
		assert.Nil(t, filter.values["Name"])
	})
}

func TestCreateActiveFilter(t *testing.T) {
	spec := OnlyActive()
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 2)
}

func TestQuerySpecification_ApplyFilters(t *testing.T) {
	assertFilters := func(querySpec *QueryBuilder, length int) {
		spec := querySpec.ToSpecification()
		builder := spec(BuildContext(), expression.NewBuilder())
		expression, _ := builder.Build()
		assert.Len(t, expression.Names(), length)
	}
	t.Run("Without filters", func(t *testing.T) {
		spec := new(QueryBuilder)
		assertFilters(spec, 1)
	})
	t.Run("With one filter", func(t *testing.T) {
		spec := new(QueryBuilder)
		spec.AddFilter("name", "test")
		assertFilters(spec, 2)
	})
	t.Run("With two filter", func(t *testing.T) {
		spec := new(QueryBuilder)
		spec.AddFilter("name", "test")
		spec.AddFilter("active", true)
		assertFilters(spec, 3)
	})
	t.Run("With three filter", func(t *testing.T) {
		spec := new(QueryBuilder)
		spec.AddFilter("name", "test")
		spec.AddFilter("active", true)
		spec.AddFilter("gender", "M")
		assertFilters(spec, 4)
	})
}

func TestCreateMarriageFilter(t *testing.T) {
	spec := OnlyMarriage()
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 3)
}

func TestOnlyLegalMembers(t *testing.T) {
	assert.True(t, OnlyLegalMembers()(BuildAdult()))
	assert.False(t, OnlyLegalMembers()(BuildChildren()))
}

func TestOnlyByClassification(t *testing.T) {
	assert.True(t, OnlyByClassification(classification.YOUNG)(BuildAdult()))
	assert.False(t, OnlyByClassification(classification.ADULT)(BuildChildren()))
}

func TestLastMarriages(t *testing.T) {
	spec := LastMarriages(time.Now(), time.Now())
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 4)
}

func TestLastBirths(t *testing.T) {
	spec := LastBirths(time.Now(), time.Now())
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 3)
}

func TestBirthDay(t *testing.T) {
	spec := WithBirthday(time.Now())
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 3)
}

func BuildChildren() *domain.Member {
	return &domain.Member{
		Person: domain.Person{
			BirthDate: time.Now(),
		},
	}
}

func BuildAdult() *domain.Member {
	return &domain.Member{
		Person: domain.Person{
			BirthDate: time.Now().AddDate(-20, 0, 0),
		},
	}
}

func BuildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID: "church_id_test",
		},
	})
}
