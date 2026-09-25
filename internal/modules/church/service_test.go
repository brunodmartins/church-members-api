package church

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/role"
	mock_church "github.com/brunodmartins/church-members-api/internal/modules/church/mock"
	mock_member "github.com/brunodmartins/church-members-api/internal/modules/member/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
)

func TestChurchService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_church.NewMockRepository(ctrl)
	service := NewService(nil, repo)

	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().List(gomock.Any()).Return([]*domain.Church{buildChurch("")}, nil)
		result, err := service.List(nil)
		assert.Nil(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().List(gomock.Any()).Return(nil, genericError)
		_, err := service.List(nil)
		assert.NotNil(t, err)
	})
}

func TestChurchService_GetChurch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_church.NewMockRepository(ctrl)
	service := NewService(nil, repo)
	const id = "xxx"
	var church = buildChurch(id)
	var ctx = context.WithValue(context.TODO(), "church", church)
	t.Run("Success", func(t *testing.T) {
		repo.EXPECT().GetByID(gomock.Any(), id).Return(church, nil)
		result, err := service.GetChurch(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, id, result.ID)
	})

	t.Run("Fail", func(t *testing.T) {
		repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, genericError)
		_, err := service.GetChurch(ctx, id)
		assert.NotNil(t, err)
	})
}

func TestChurchService_GetChurchByAbbreviation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_church.NewMockRepository(ctrl)
	service := NewService(nil, repo)
	const abbreviation = "TC"
	const id = "xxx"

	t.Run("Success - Church Found", func(t *testing.T) {
		expected := buildChurch(id)
		expected.Abbreviation = abbreviation
		churches := []*domain.Church{expected}

		repo.EXPECT().List(gomock.Any()).Return(churches, nil)

		result, err := service.GetChurchByAbbreviation(nil, abbreviation)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error - Repository Error", func(t *testing.T) {
		repo.EXPECT().List(gomock.Any()).Return(nil, errors.New("database error"))

		result, err := service.GetChurchByAbbreviation(nil, abbreviation)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - Church Not Found", func(t *testing.T) {
		expected := buildChurch(id)
		expected.Abbreviation = abbreviation
		churches := []*domain.Church{expected}

		repo.EXPECT().List(gomock.Any()).Return(churches, nil)

		result, err := service.GetChurchByAbbreviation(nil, "OTHER")
		assert.Error(t, err)
		assert.Nil(t, result)

		// Check if it's an API error with correct status
		apiErr, ok := err.(apierrors.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})
}

func TestChurchService_GetStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_church.NewMockRepository(ctrl)
	memberService := mock_member.NewMockService(ctrl)
	service := NewService(memberService, repo)
	const id = "xxx"
	var church = buildChurch(id)
	var ctx = context.WithValue(context.TODO(), "church", church)
	t.Run("Forbidden Access", func(t *testing.T) {
		stats, err := service.GetStatistics(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, stats)
		apiErr, ok := err.(apierrors.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, apiErr.StatusCode())
	})

	t.Run("Member Service Error", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(ctx, gomock.Any()).Return(nil, errors.New("service error"))
		stats, err := service.GetStatistics(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, stats)
	})

	t.Run("Success", func(t *testing.T) {
		members := []*domain.Member{
			{Person: &domain.Person{BirthDate: time.Now(), Gender: "M"}},
			{Person: &domain.Person{BirthDate: time.Now(), Gender: "F"}},
		}
		memberService.EXPECT().SearchMembers(ctx, gomock.Any()).Return(members, nil)

		stats, err := service.GetStatistics(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 2, stats.TotalMembers)
		assert.Len(t, stats.AgeDistribution, 2)
		assert.Equal(t, 1, stats.TotalMembersByGender["M"])
		assert.Equal(t, 1, stats.TotalMembersByGender["F"])
	})
}

func TestChurchService_UpdateChurch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_church.NewMockRepository(ctrl)
	service := NewService(nil, repo)
	const id = "xxx"

	t.Run("Success", func(t *testing.T) {
		currentChurch := buildChurch(id)
		ctx := context.WithValue(context.TODO(), "user", &domain.User{
			Role: role.ADMIN,
			Church: &domain.Church{
				ID: id,
			},
		})
		update := &domain.Church{ID: id, Name: "updated church", Language: "en-us", Email: "church@example.com", Logo: "logo.png"}
		repo.EXPECT().GetByID(gomock.Eq(ctx), id).Return(currentChurch, nil)
		repo.EXPECT().Update(gomock.Eq(ctx), gomock.Eq(&domain.Church{
			ID:           id,
			Name:         update.Name,
			Abbreviation: currentChurch.Abbreviation,
			Language:     update.Language,
			Email:        update.Email,
			Logo:         update.Logo,
		})).Return(nil)
		assert.NoError(t, service.UpdateChurch(ctx, update))
	})

	t.Run("Forbidden for non-admin", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), "user", &domain.User{
			Role: role.USER,
			Church: &domain.Church{
				ID: id,
			},
		})
		err := service.UpdateChurch(ctx, &domain.Church{ID: id})
		assert.Error(t, err)
		assert.Equal(t, http.StatusForbidden, err.(apierrors.Error).StatusCode())
	})

	t.Run("Get church error", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), "user", &domain.User{
			Role: role.ADMIN,
			Church: &domain.Church{
				ID: id,
			},
		})
		repo.EXPECT().GetByID(gomock.Eq(ctx), id).Return(nil, genericError)
		assert.Error(t, service.UpdateChurch(ctx, &domain.Church{ID: id}))
	})

	t.Run("Forbidden for admin from other church", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), "user", &domain.User{
			Role: role.ADMIN,
			Church: &domain.Church{
				ID: "other-church",
			},
		})
		err := service.UpdateChurch(ctx, &domain.Church{ID: id})
		assert.Error(t, err)
		assert.Equal(t, http.StatusForbidden, err.(apierrors.Error).StatusCode())
	})
}
