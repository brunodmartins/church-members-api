package participant_test

import (
	"context"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/participant"
	mock_participant "github.com/brunodmartins/church-members-api/internal/modules/participant/mock"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSearchParticipants(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_participant.NewMockRepository(ctrl)
	service := participant.NewService(repo)

	spec := wrapper.QuerySpecification(nil)

	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(BuildParticipants(2), nil)
		participants, err := service.SearchParticipant(BuildContext(), wrapper.QuerySpecification(nil))
		assert.Nil(t, err)
		assert.Len(t, participants, 2)
	})

	t.Run("Success with post specification", func(t *testing.T) {
		parts := BuildParticipants(2)
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(parts, nil)
		filter := participant.Specification(func(p *domain.Participant) bool { return p.ID == parts[0].ID })
		participants, err := service.SearchParticipant(BuildContext(), wrapper.QuerySpecification(nil), filter)
		assert.Nil(t, err)
		assert.Len(t, participants, 1)
		assert.Equal(t, parts[0].ID, participants[0].ID)
	})

	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		_, err := service.SearchParticipant(BuildContext(), wrapper.QuerySpecification(nil))
		assert.NotNil(t, err)
	})
}

func TestGetParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_participant.NewMockRepository(ctrl)
	service := participant.NewService(repo)

	id := domain.NewID()
	part := buildParticipant(id)

	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(part, nil)
		found, err := service.GetParticipant(BuildContext(), id)
		assert.Nil(t, err)
		assert.Equal(t, id, found.ID)
	})

	t.Run("Fail - repo error", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(nil, genericError)
		_, err := service.GetParticipant(BuildContext(), id)
		assert.NotNil(t, err)
	})

	t.Run("Fail - invalid id", func(t *testing.T) {
		_, err := service.GetParticipant(BuildContext(), "")
		assert.NotNil(t, err)
	})
}

func TestCreateUpdateDeleteParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_participant.NewMockRepository(ctrl)
	service := participant.NewService(repo)

	t.Run("Create Success", func(t *testing.T) {
		p := buildParticipant("")
		repo.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(p)).DoAndReturn(func(ctx context.Context, participant *domain.Participant) error {
			participant.ID = domain.NewID()
			return nil
		})
		id, err := service.CreateParticipant(BuildContext(), p)
		assert.Nil(t, err)
		assert.NotEmpty(t, id)
		assert.Equal(t, "church_id_test", p.ChurchID)
	})

	t.Run("Create Fail", func(t *testing.T) {
		p := buildParticipant("")
		repo.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(p)).Return(genericError)
		_, err := service.CreateParticipant(BuildContext(), p)
		assert.NotNil(t, err)
	})

	t.Run("Update Success", func(t *testing.T) {
		p := buildParticipant(domain.NewID())
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(p.ID)).Return(buildParticipant(p.ID), nil)
		repo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(p)).Return(nil)
		err := service.UpdateParticipant(BuildContext(), p)
		assert.Nil(t, err)
	})

	t.Run("Update Fail", func(t *testing.T) {
		p := buildParticipant(domain.NewID())
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(p.ID)).Return(buildParticipant(p.ID), nil)
		repo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(p)).Return(genericError)
		err := service.UpdateParticipant(BuildContext(), p)
		assert.NotNil(t, err)
	})

	t.Run("Delete Success", func(t *testing.T) {
		id := domain.NewID()
		repo.EXPECT().Delete(gomock.Any(), gomock.Eq(id)).Return(nil)
		err := service.DeleteParticipant(BuildContext(), id)
		assert.Nil(t, err)
	})

	t.Run("Delete Fail", func(t *testing.T) {
		id := domain.NewID()
		repo.EXPECT().Delete(gomock.Any(), gomock.Eq(id)).Return(genericError)
		err := service.DeleteParticipant(BuildContext(), id)
		assert.NotNil(t, err)
	})
}
