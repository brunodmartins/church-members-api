package participant_test

import (
	"context"
	"testing"
	"time"

	"github.com/brunodmartins/church-members-api/internal/modules/participant"
	mock_participant "github.com/brunodmartins/church-members-api/internal/modules/participant/mock"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"go.uber.org/mock/gomock"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/stretchr/testify/assert"
)

func TestListAllParticipants(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_participant.NewMockRepository(ctrl)
	service := participant.NewService(repo)
	spec := wrapper.QuerySpecification(nil)

	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(BuildParticipants(2), nil)
		parts, err := service.SearchParticipant(BuildContext(), wrapper.QuerySpecification(nil))
		assert.Nil(t, err)
		assert.Len(t, parts, 2)
	})

	t.Run("Success with post specification", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(BuildParticipants(2), nil)
		// post specification that matches the BuildParticipants name
		postSpec := participant.Specification(func(p *domain.Participant) bool { return p.Name == "First Last" })
		parts, err := service.SearchParticipant(BuildContext(), wrapper.QuerySpecification(nil), postSpec)
		assert.Nil(t, err)
		assert.Len(t, parts, 2)
	})

	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().FindAll(gomock.Any(), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		_, err := service.SearchParticipant(BuildContext(), wrapper.QuerySpecification(nil))
		assert.NotNil(t, err)
	})
}

func TestFindParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_participant.NewMockRepository(ctrl)
	service := participant.NewService(repo)

	id := domain.NewID()
	part := buildParticipant(id)

	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(part, nil)
		found, err := service.GetParticipant(BuildContext(), id)
		assert.Equal(t, id, found.ID)
		assert.Nil(t, err)
	})

	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(nil, genericError)
		_, err := service.GetParticipant(BuildContext(), id)
		assert.NotNil(t, err)
	})

	t.Run("Fail - Invalid ID", func(t *testing.T) {
		_, err := service.GetParticipant(BuildContext(), "")
		assert.NotNil(t, err)
	})
}

func TestCreateParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_participant.NewMockRepository(ctrl)
	service := participant.NewService(repo)

	t.Run("Success", func(t *testing.T) {
		p := buildParticipant("")
		repo.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(p)).DoAndReturn(func(ctx context.Context, participant *domain.Participant) error {
			participant.ID = domain.NewID()
			return nil
		})
		id, err := service.CreateParticipant(BuildContext(), p)
		assert.Nil(t, err)
		assert.NotEmpty(t, p.ID)
		assert.NotEmpty(t, id)
	})

	t.Run("Fail", func(t *testing.T) {
		p := buildParticipant("")
		repo.EXPECT().Insert(gomock.Any(), gomock.AssignableToTypeOf(p)).Return(genericError)
		_, err := service.CreateParticipant(BuildContext(), p)
		assert.NotNil(t, err)
	})
}

func TestUpdateParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_participant.NewMockRepository(ctrl)
	service := participant.NewService(repo)

	t.Run("Success", func(t *testing.T) {
		id := domain.NewID()
		existing := buildParticipant(id)
		update := &domain.Participant{ID: id, Name: "New Name", Filiation: "New fil", BirthDate: existing.BirthDate, Observation: "Obs2", CellPhone: "1111"}
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(existing, nil)
		repo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(existing)).Return(nil)
		assert.NoError(t, service.UpdateParticipant(BuildContext(), update))
	})

	t.Run("Fail - Update", func(t *testing.T) {
		id := domain.NewID()
		existing := buildParticipant(id)
		update := &domain.Participant{ID: id}
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(existing, nil)
		repo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(existing)).Return(genericError)
		assert.Error(t, service.UpdateParticipant(BuildContext(), update))
	})

	t.Run("Fail - GetParticipant", func(t *testing.T) {
		id := domain.NewID()
		update := &domain.Participant{ID: id}
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(nil, genericError)
		assert.Error(t, service.UpdateParticipant(BuildContext(), update))
	})
}

func TestRetireParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_participant.NewMockRepository(ctrl)
	service := participant.NewService(repo)

	id := domain.NewID()
	reason := "testing"
	date := time.Now()
	existing := buildParticipant(id)

	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(existing, nil)
		repo.EXPECT().RetireParticipant(gomock.Any(), gomock.AssignableToTypeOf(existing)).Return(nil)
		assert.NoError(t, service.RetireParticipant(BuildContext(), id, reason, date))
	})

	t.Run("Fail - Retire", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(existing, nil)
		repo.EXPECT().RetireParticipant(gomock.Any(), gomock.AssignableToTypeOf(existing)).Return(genericError)
		assert.Error(t, service.RetireParticipant(BuildContext(), id, reason, date))
	})

	t.Run("Fail - GetParticipant", func(t *testing.T) {
		repo.EXPECT().FindByID(gomock.Any(), gomock.Eq(id)).Return(nil, genericError)
		assert.Error(t, service.RetireParticipant(BuildContext(), id, reason, date))
	})
}
