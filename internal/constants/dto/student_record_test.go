package dto

import (
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/stretchr/testify/assert"
)

func TestStudentItemRoundTrip(t *testing.T) {
	freshStudent := func() *domain.Student {
		return &domain.Student{
			ID:       domain.NewID(),
			ChurchID: "church-id",
			Name:     "Student Name",
			Gender:   "F",
		}
	}

	t.Run("with member", func(t *testing.T) {
		student := freshStudent()
		student.Member = &domain.Member{
			ID: domain.NewID(),
			Person: &domain.Person{
				FirstName: "Member",
				LastName:  "Name",
				Gender:    "F",
				Contact:   &domain.Contact{},
				Address:   &domain.Address{},
			},
			Religion: &domain.Religion{},
		}

		item := NewStudentItem(student)
		got := item.ToStudent()

		assert.Equal(t, student.ID, got.ID)
		assert.Equal(t, student.ChurchID, item.ChurchID)
		assert.Equal(t, student.ChurchID, got.ChurchID)
		assert.Equal(t, student.Name, got.Name)
		assert.Equal(t, student.Gender, got.Gender)
		assert.Equal(t, student.Member.ID, got.Member.ID)
		assert.Nil(t, got.Participant)
	})

	t.Run("with participant", func(t *testing.T) {
		student := freshStudent()
		student.Participant = &domain.Participant{
			ID:     domain.NewID(),
			Name:   "Participant Name",
			Gender: "M",
		}

		item := NewStudentItem(student)
		got := item.ToStudent()

		assert.Equal(t, student.ID, got.ID)
		assert.Equal(t, student.ChurchID, item.ChurchID)
		assert.Equal(t, student.ChurchID, got.ChurchID)
		assert.Equal(t, student.Name, got.Name)
		assert.Equal(t, student.Gender, got.Gender)
		assert.Equal(t, student.Participant.ID, got.Participant.ID)
		assert.Nil(t, got.Member)
	})
}
