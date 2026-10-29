package student_test

import (
	"context"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	mock_member "github.com/brunodmartins/church-members-api/internal/modules/member/mock"
	mock_participant "github.com/brunodmartins/church-members-api/internal/modules/participant/mock"
	"github.com/brunodmartins/church-members-api/internal/modules/student"
	mock_student "github.com/brunodmartins/church-members-api/internal/modules/student/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := mock_student.NewMockRepository(ctrl)
	memberService := mock_member.NewMockService(ctrl)
	participantService := mock_participant.NewMockService(ctrl)
	service := student.NewService(repository, memberService, participantService)
	ctx := context.WithValue(context.Background(), "church_id", domain.NewID())

	t.Run("creates a fresh student", func(t *testing.T) {
		createdStudent := &domain.Student{Name: "Jane Doe", Gender: "F"}
		repository.EXPECT().Insert(gomock.Any(), createdStudent).DoAndReturn(func(_ context.Context, student *domain.Student) error {
			student.ID = domain.NewID()
			return nil
		})

		id, err := service.CreateStudent(ctx, createdStudent)

		assert.NoError(t, err)
		assert.Equal(t, createdStudent.ID, id)
		assert.Equal(t, ctx.Value("church_id"), createdStudent.ChurchID)
	})

	t.Run("creates a student from a member", func(t *testing.T) {
		memberID := domain.NewID()
		existingMember := &domain.Member{
			ID: memberID,
			Person: &domain.Person{
				FirstName: "Jane",
				LastName:  "Member",
				Gender:    "F",
			},
		}
		createdStudent := &domain.Student{Member: &domain.Member{ID: memberID}}
		memberService.EXPECT().GetMember(gomock.Any(), memberID).Return(existingMember, nil)
		repository.EXPECT().Insert(gomock.Any(), createdStudent).DoAndReturn(func(_ context.Context, student *domain.Student) error {
			assert.Equal(t, "Jane Member", student.Name)
			assert.Equal(t, "F", student.Gender)
			assert.Same(t, existingMember, student.Member)
			student.ID = domain.NewID()
			return nil
		})

		_, err := service.CreateStudent(ctx, createdStudent)

		assert.NoError(t, err)
	})

	t.Run("creates a student from a participant", func(t *testing.T) {
		participantID := domain.NewID()
		existingParticipant := &domain.Participant{ID: participantID, Name: "John Participant", Gender: "M"}
		createdStudent := &domain.Student{Participant: &domain.Participant{ID: participantID}}
		participantService.EXPECT().GetParticipant(gomock.Any(), participantID).Return(existingParticipant, nil)
		repository.EXPECT().Insert(gomock.Any(), createdStudent).DoAndReturn(func(_ context.Context, student *domain.Student) error {
			assert.Equal(t, "John Participant", student.Name)
			assert.Equal(t, "M", student.Gender)
			assert.Same(t, existingParticipant, student.Participant)
			student.ID = domain.NewID()
			return nil
		})

		_, err := service.CreateStudent(ctx, createdStudent)

		assert.NoError(t, err)
	})

	t.Run("rejects multiple sources", func(t *testing.T) {
		errStudent := &domain.Student{
			Member:      &domain.Member{ID: domain.NewID()},
			Participant: &domain.Participant{ID: domain.NewID()},
		}

		_, err := service.CreateStudent(ctx, errStudent)

		assert.Error(t, err)
	})
}
