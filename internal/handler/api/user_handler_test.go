package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	mock_user "github.com/brunodmartins/church-members-api/internal/modules/user/mock"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func TestUserHandler_PostUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_user.NewMockService(ctrl)
	NewUserHandler(service).SetUpRoutes(app)

	t.Run("Success", func(t *testing.T) {
		service.EXPECT().SaveUser(gomock.Any(), gomock.AssignableToTypeOf(new(domain.User))).Return(nil)
		body := getMock("create_user_valid.json")
		runTest(app, buildPost("/users", body)).assertStatus(t, http.StatusCreated)
	})
	t.Run("Fail - Service error - 500", func(t *testing.T) {
		service.EXPECT().SaveUser(gomock.Any(), gomock.AssignableToTypeOf(new(domain.User))).Return(genericError)
		body := getMock("create_user_valid.json")
		runTest(app, buildPost("/users", body)).assertStatus(t, http.StatusInternalServerError)
	})
	t.Run("Fail - empty body - 400", func(t *testing.T) {
		runTest(app, buildPost("/users", emptyJson)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - invalid role - 400", func(t *testing.T) {
		body := getMock("create_user_invalid_role.json")
		runTest(app, buildPost("/users", body)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - invalid email - 400", func(t *testing.T) {
		body := getMock("create_user_invalid_email.json")
		runTest(app, buildPost("/users", body)).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Fail - invalid password - 400", func(t *testing.T) {
		body := getMock("create_user_invalid_password.json")
		runTest(app, buildPost("/users", body)).assertStatus(t, http.StatusBadRequest)
	})
}

func TestUserHandler_ConfirmUserEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_user.NewMockService(ctrl)
	NewUserHandler(service).SetUpPublicRoutes(app)

	userID := domain.NewID()
	churchID := domain.NewID()
	token := "token"

	ctx := context.WithValue(context.Background(), "church_id", churchID)
	t.Run("Given a valid URL, when perform a GET operation, then return 200 OK", func(t *testing.T) {
		service.EXPECT().ConfirmEmail(gomock.Eq(ctx), gomock.Eq(userID), gomock.Eq(token)).Return(nil)
		runTest(app, buildGet(fmt.Sprintf("/users/%s/confirm?church=%s&token=%s", userID, churchID, token))).assertStatus(t, http.StatusOK)
	})
	t.Run("Given a valid URL, when perform a GET operation, then return 500 error", func(t *testing.T) {
		service.EXPECT().ConfirmEmail(gomock.Eq(ctx), gomock.Eq(userID), gomock.Eq(token)).Return(errors.New("generic error"))
		runTest(app, buildGet(fmt.Sprintf("/users/%s/confirm?church=%s&token=%s", userID, churchID, token))).assertStatus(t, http.StatusInternalServerError)
	})
	t.Run("Given a invalid URL (without token), when perform a GET operation, then return 400 bad request", func(t *testing.T) {
		runTest(app, buildGet(fmt.Sprintf("/users/%s/confirm?church=%s", userID, churchID))).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Given a invalid URL (without church), when perform a GET operation, then return 400 bad request", func(t *testing.T) {
		runTest(app, buildGet(fmt.Sprintf("/users/%s/confirm?token=%s", userID, token))).assertStatus(t, http.StatusBadRequest)
	})
	t.Run("Given a invalid URL (without church and token), when perform a GET operation, then return 400 bad request", func(t *testing.T) {
		runTest(app, buildGet(fmt.Sprintf("/users/%s/confirm", userID))).assertStatus(t, http.StatusBadRequest)
	})
}
