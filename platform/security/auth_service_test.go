package security

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/services/email"
	mock_email "github.com/brunodmartins/church-members-api/internal/services/email/mock"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	mock_church "github.com/brunodmartins/church-members-api/internal/modules/church/mock"
	mock_user "github.com/brunodmartins/church-members-api/internal/modules/user/mock"
	"github.com/brunodmartins/church-members-api/platform/crypto"
	"github.com/spf13/viper"

	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mock_user.NewMockService(ctrl)
	churchService := mock_church.NewMockService(ctrl)
	emailService := mock_email.NewMockService(ctrl)
	service := NewAuthService(userService, churchService, emailService)
	church := buildChurch(domain.NewID())
	testUser := buildUser("", string(crypto.EncryptPassword(password)))
	t.Run("Given a valid user, when try to authenticate, then succeeds", func(t *testing.T) {
		userService.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(testUser, nil)
		token, err := service.GenerateToken(church, userName, password)
		assert.NotEmpty(t, token)
		assert.Nil(t, err)
	})

	t.Run("Given a valid user, when try to authenticate, then fails due to password not match", func(t *testing.T) {
		userService.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(testUser, nil)
		token, err := service.GenerateToken(church, userName, password+"123")
		assert.Empty(t, token)
		assert.Error(t, err)
	})
	t.Run("Given a valid user, when try to authenticate, then fails due to error on look up for the user", func(t *testing.T) {
		userService.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(nil, genericError)
		token, err := service.GenerateToken(church, userName, password)
		assert.Empty(t, token)
		assert.Error(t, err)
	})
	t.Run("Given a valid user, when try to authenticate, then fails due to user not have confirmed the email, and send email", func(t *testing.T) {
		unconfirmedUser := buildUser("", string(crypto.EncryptPassword(password)))
		unconfirmedUser.ConfirmedEmail = false
		userService.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(unconfirmedUser, nil)
		emailService.EXPECT().SendTemplateEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		token, err := service.GenerateToken(church, userName, password)
		assert.Empty(t, token)
		assert.Error(t, err)
	})
	t.Run("Given a valid user, when try to authenticate, then fails due to user not have confirmed the email, and send email fails", func(t *testing.T) {
		unconfirmedUser := buildUser("", string(crypto.EncryptPassword(password)))
		unconfirmedUser.ConfirmedEmail = false
		userService.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(unconfirmedUser, nil)
		emailService.EXPECT().SendTemplateEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(genericError)
		token, err := service.GenerateToken(church, userName, password)
		assert.Empty(t, token)
		assert.Error(t, err)
	})
}

func TestAuthService_IsValidToken(t *testing.T) {
	viper.Set("security.token.expiration", 1)
	valid, _ := GetClaim("")
	assert.False(t, valid)
	valid, _ = GetClaim(buildToken())
	assert.True(t, valid)
}

func TestAddClaimToContext(t *testing.T) {
	viper.Set("security.token.expiration", 1)
	ctx := context.Background()
	_, claim := GetClaim(buildToken())
	assert.Nil(t, ctx.Value("user"))
	ctx = AddClaimToContext(claim, ctx)
	assert.NotNil(t, ctx.Value("user"))
}

func TestService_SendConfirmationEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mock_user.NewMockService(ctrl)
	churchService := mock_church.NewMockService(ctrl)
	emailService := mock_email.NewMockService(ctrl)
	service := NewAuthService(userService, churchService, emailService)
	user := buildUser("id", "")
	t.Run("Given a valid user, when send the confirmation email, then send confirmation email successfully", func(t *testing.T) {
		emailService.EXPECT().SendTemplateEmail(gomock.Eq(email.ConfirmEmailTemplate), gomock.AssignableToTypeOf(email.ConfirmEmailTemplateDTO{}), gomock.Any(), gomock.Eq(user.Email)).Return(nil)
		assert.NoError(t, service.SendConfirmEmail(BuildContext(), user))
	})
	t.Run("Given a valid user, when send the confirmation email, then send confirmation email fails", func(t *testing.T) {
		emailService.EXPECT().SendTemplateEmail(gomock.Eq(email.ConfirmEmailTemplate), gomock.AssignableToTypeOf(email.ConfirmEmailTemplateDTO{}), gomock.Any(), gomock.Eq(user.Email)).Return(genericError)
		assert.Error(t, service.SendConfirmEmail(BuildContext(), user))
	})
}

func TestUserService_ConfirmEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mock_user.NewMockService(ctrl)
	churchService := mock_church.NewMockService(ctrl)
	emailService := mock_email.NewMockService(ctrl)
	service := NewAuthService(userService, churchService, emailService)

	t.Run("Given a valid user and token, when the confirm operation is call, then confirm correctly", func(t *testing.T) {
		ctx := BuildContext()
		user := buildUser(domain.NewID(), "")
		user.ConfirmedEmail = false
		userService.EXPECT().FindUser(gomock.Eq(ctx), gomock.Eq(userName)).Return(user, nil)
		userService.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(user)).Return(nil)
		assert.Nil(t, service.ConfirmEmail(ctx, userName))
		assert.True(t, user.ConfirmedEmail)
	})
	t.Run("Given a valid user and token, when the confirm operation is call, then fails the operation due to repository error", func(t *testing.T) {
		ctx := BuildContext()
		user := buildUser(domain.NewID(), "")
		user.ConfirmedEmail = false
		userService.EXPECT().FindUser(gomock.Eq(ctx), gomock.Eq(userName)).Return(user, nil)
		userService.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(user)).Return(errors.New("generic error"))
		assert.Error(t, service.ConfirmEmail(ctx, userName))
	})
	t.Run("Given a user that has already confirmed, when the confirm operation is call, then fails the operation due to already be confirmed", func(t *testing.T) {
		ctx := BuildContext()
		user := buildUser(domain.NewID(), "")
		user.ConfirmedEmail = true
		userService.EXPECT().FindUser(gomock.Eq(ctx), gomock.Eq(userName)).Return(user, nil)
		err := service.ConfirmEmail(ctx, userName)
		assert.Error(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, err.(apierrors.Error).StatusCode())
	})
	t.Run("Given a nonexistent user, when the confirm operation is call, then fails the operation due to not found the user", func(t *testing.T) {
		ctx := BuildContext()
		user := buildUser(domain.NewID(), "")
		user.ConfirmedEmail = true
		userService.EXPECT().FindUser(gomock.Eq(ctx), gomock.Eq(userName)).Return(nil, apierrors.NewApiError("User not found", http.StatusNotFound))
		err := service.ConfirmEmail(ctx, userName)
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
	})

}

func TestAuthService_IdentifyChurch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mock_user.NewMockService(ctrl)
	churchService := mock_church.NewMockService(ctrl)
	emailService := mock_email.NewMockService(ctrl)
	service := NewAuthService(userService, churchService, emailService)
	church := buildChurch(domain.NewID())
	ctx := context.TODO()

	t.Run("Invalid arguments", func(t *testing.T) {
		_, err := service.IdentifyChurch(ctx, "", "")
		assert.Equal(t, apierrors.NewApiError("Church ID or abbreviation must be provided", http.StatusBadRequest), err)
	})
	t.Run("Success on Church Abbreviation", func(t *testing.T) {
		churchService.EXPECT().GetChurchByAbbreviation(gomock.Eq(ctx), church.Abbreviation).Return(church, nil)
		result, err := service.IdentifyChurch(ctx, church.Abbreviation, "")
		assert.Nil(t, err)
		assert.Equal(t, church, result)
	})
	t.Run("Success on Church ID", func(t *testing.T) {
		churchService.EXPECT().GetChurch(gomock.Eq(ctx), church.ID).Return(church, nil)
		result, err := service.IdentifyChurch(ctx, "", church.ID)
		assert.Nil(t, err)
		assert.Equal(t, church, result)
	})
	t.Run("Fails on Church Abbreviation", func(t *testing.T) {
		churchService.EXPECT().GetChurchByAbbreviation(gomock.Eq(ctx), church.Abbreviation).Return(nil, errors.New("error"))
		result, err := service.IdentifyChurch(ctx, church.Abbreviation, "")
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
	t.Run("Fails on Church ID", func(t *testing.T) {
		churchService.EXPECT().GetChurch(gomock.Eq(ctx), church.ID).Return(nil, errors.New("error"))
		result, err := service.IdentifyChurch(ctx, "", church.ID)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func Test_buildConfirmationLink(t *testing.T) {
	viper.Set("email.confirm.url", "http://localhost")
	accessToken := "my-token"
	expected := fmt.Sprintf("http://localhost/users/confirm?accessToken=%s", accessToken)
	assert.Equal(t, expected, buildConfirmationLink(accessToken))
}

func TestAuthService_GetRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userService := mock_user.NewMockService(ctrl)
	churchService := mock_church.NewMockService(ctrl)
	emailService := mock_email.NewMockService(ctrl)
	service := NewAuthService(userService, churchService, emailService)

	ctx := BuildContext()
	viper.Set("security.token.expiration", 1)
	jwtToken := GenerateJWTToken(buildUser("test_user", "id", "admin"))
	roles := []string{"admin"}

	t.Run("Success on GetRoles", func(t *testing.T) {
		result, err := service.GetRoles(ctx, jwtToken)
		assert.Nil(t, err)
		assert.Equal(t, roles, result)
	})
	t.Run("Fail on GetRoles", func(t *testing.T) {
		_, err := service.GetRoles(ctx, "")
		assert.NotNil(t, err)
	})
}

func BuildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID: "church_id_test",
		},
	})
}

func buildChurch(id string) *domain.Church {
	return &domain.Church{
		ID:           id,
		Name:         "test church",
		Abbreviation: "tc",
		Language:     "pt-br",
		Email:        "test@test.com",
	}
}
