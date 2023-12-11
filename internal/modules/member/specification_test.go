package member

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/classification"

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
	queryExpression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, queryExpression.Names(), 2)
}

func TestQuerySpecification_ApplyFilters(t *testing.T) {
	t.Run("Query without filters", func(t *testing.T) {
		querySpec := new(QueryBuilder)
		spec := querySpec.ToSpecification()
		builder := spec(BuildContext(), expression.NewBuilder())
		queryExpression, err := builder.Build()
		assert.Nil(t, err)
		assert.NotNil(t, queryExpression.KeyCondition())
		assert.Nil(t, queryExpression.Condition())
		assert.Nil(t, queryExpression.Filter())
		assert.Len(t, queryExpression.Names(), 1)
	})
	t.Run("Query with name filter", func(t *testing.T) {
		querySpec := new(QueryBuilder)
		querySpec.AddFilter("name", "test")
		spec := querySpec.ToSpecification()
		builder := spec(BuildContext(), expression.NewBuilder())
		queryExpression, err := builder.Build()
		assert.Nil(t, err)
		assert.NotNil(t, queryExpression.KeyCondition())
		assert.Nil(t, queryExpression.Condition())
		assert.Nil(t, queryExpression.Filter())
		assert.Len(t, queryExpression.Names(), 2)
	})
	t.Run("Query with name and active filter", func(t *testing.T) {
		querySpec := new(QueryBuilder)
		querySpec.AddFilter("name", "test")
		querySpec.AddFilter("active", true)
		spec := querySpec.ToSpecification()
		builder := spec(BuildContext(), expression.NewBuilder())
		queryExpression, err := builder.Build()
		assert.Nil(t, err)
		assert.NotNil(t, queryExpression.KeyCondition())
		assert.Nil(t, queryExpression.Condition())
		assert.NotNil(t, queryExpression.Filter())
		assert.Len(t, queryExpression.Names(), 3)
	})
	t.Run("With three filter", func(t *testing.T) {
		querySpec := new(QueryBuilder)
		querySpec.AddFilter("name", "test")
		querySpec.AddFilter("active", true)
		querySpec.AddFilter("gender", "M")
		spec := querySpec.ToSpecification()
		builder := spec(BuildContext(), expression.NewBuilder())
		queryExpression, err := builder.Build()
		assert.Nil(t, err)
		assert.NotNil(t, queryExpression.KeyCondition())
		assert.Nil(t, queryExpression.Condition())
		assert.NotNil(t, queryExpression.Filter())
		assert.Len(t, queryExpression.Names(), 4)
	})
}

func TestCreateMarriageFilter(t *testing.T) {
	spec := OnlyMarriage()
	builder := spec(BuildContext(), expression.NewBuilder())
	queryExpression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, queryExpression.Names(), 3)
}

func TestOnlyLegalMembers(t *testing.T) {
	assert.True(t, OnlyLegalMembers()(BuildAdult()))
	assert.False(t, OnlyLegalMembers()(BuildChildren()))
}

func TestOnlyMembershipEndCurrentYear(t *testing.T) {
	t.Run("Member active", func(t *testing.T) {
		member := BuildAdult()
		member.Active = true
		member.MembershipEndDate = nil
		assert.False(t, OnlyMembershipEndCurrentYear()(member))
	})
	t.Run("Member Inactive current year", func(t *testing.T) {
		member := BuildAdult()
		member.Active = false
		now := time.Now()
		member.MembershipEndDate = &now
		assert.True(t, OnlyMembershipEndCurrentYear()(member))
	})
	t.Run("Member Inactive last year", func(t *testing.T) {
		member := BuildAdult()
		member.Active = false
		now := time.Now().AddDate(-1, 0, 0)
		member.MembershipEndDate = &now
		assert.False(t, OnlyMembershipEndCurrentYear()(member))
	})
}

func TestOnlyByClassification(t *testing.T) {
	assert.True(t, OnlyByClassification(classification.YOUNG)(BuildAdult()))
	assert.False(t, OnlyByClassification(classification.ADULT)(BuildChildren()))
}

func TestLastMarriages(t *testing.T) {
	spec := LastMarriages(time.Now(), time.Now())
	builder := spec(BuildContext(), expression.NewBuilder())
	queryExpression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, queryExpression.Names(), 4)
}

func TestLastBirths(t *testing.T) {
	spec := LastBirths(time.Now(), time.Now())
	builder := spec(BuildContext(), expression.NewBuilder())
	queryExpression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, queryExpression.Names(), 3)
}

func TestBirthDay(t *testing.T) {
	spec := WithBirthday(time.Now())
	builder := spec(BuildContext(), expression.NewBuilder())
	queryExpression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, queryExpression.Names(), 3)
}

func BuildChildren() *domain.Member {
	return &domain.Member{
		Person: &domain.Person{
			BirthDate: time.Now(),
		},
	}
}

func BuildAdult() *domain.Member {
	return &domain.Member{
		Person: &domain.Person{
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
