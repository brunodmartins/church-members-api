package report

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateMarriageFilter(t *testing.T) {
	builder := expression.NewBuilder()
	spec := createMarriageFilter()
	builder = spec(builder)
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 2)
}

func TestSelectByClassification(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		var members []*domain.Member
		assert.Nil(t, selectByClassification(enum.CHILDREN, members))
	})
	t.Run("Filtering", func(t *testing.T) {
		members := []*domain.Member{
			BuildChildren(),
			BuildAdult(),
		}
		filtered := selectByClassification(enum.CHILDREN, members)
		assert.Len(t, filtered, 1)
		assert.Equal(t, enum.CHILDREN, filtered[0].Classification())

	})
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
			BirthDate: time.Now().AddDate(-20,0,0),
		},
	}
}