package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/role"
	mock_user "github.com/brunodmartins/church-members-api/internal/modules/user/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"go.uber.org/mock/gomock"
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

func TestSearchUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_user.NewMockService(ctrl)
	NewUserHandler(service).SetUpRoutes(app)

	t.Run("Success", func(t *testing.T) {
		service.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(buildUsers(), nil)
		runTest(app, buildGet("/users")).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - Service error - 500", func(t *testing.T) {
		service.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(nil, genericError)
		runTest(app, buildGet("/users")).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestGetUserByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_user.NewMockService(ctrl)
	NewUserHandler(service).SetUpRoutes(app)

	t.Run("Success", func(t *testing.T) {
		service.EXPECT().FindUser(gomock.Any(), "user1").Return(buildUsers()[0], nil)
		runTest(app, buildGet("/users/user1")).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - Service error - 500", func(t *testing.T) {
		service.EXPECT().FindUser(gomock.Any(), "user1").Return(nil, genericError)
		runTest(app, buildGet("/users/user1")).assertStatus(t, http.StatusInternalServerError)
	})
	t.Run("Fail - User not found - 404", func(t *testing.T) {
		service.EXPECT().FindUser(gomock.Any(), "user1").Return(nil, apierrors.NewApiError("User not found", 404))
		runTest(app, buildGet("/users/user1")).assertStatus(t, http.StatusNotFound)
	})
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_user.NewMockService(ctrl)
	NewUserHandler(service).SetUpRoutes(app)

	t.Run("Success", func(t *testing.T) {
		updatedUser := buildUsers()[0]
		service.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(updatedUser)).Return(nil)
		updateDTO := buildUpdateUserRequest(updatedUser)
		jsonRequest, _ := json.Marshal(updateDTO)
		runTest(app, buildPut("/users/"+updatedUser.UserName, jsonRequest)).assertStatus(t, http.StatusOK)
	})
	t.Run("Fail - Service error - 500", func(t *testing.T) {
		updatedUser := buildUsers()[0]
		service.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(updatedUser)).Return(genericError)
		updateDTO := buildUpdateUserRequest(updatedUser)
		jsonRequest, _ := json.Marshal(updateDTO)
		runTest(app, buildPut("/users/"+updatedUser.UserName, jsonRequest)).assertStatus(t, http.StatusInternalServerError)
	})
	t.Run("Fail - User not found - 404", func(t *testing.T) {
		updatedUser := buildUsers()[0]
		service.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(updatedUser)).Return(apierrors.NewApiError("User not found", 404))
		updateDTO := buildUpdateUserRequest(updatedUser)
		jsonRequest, _ := json.Marshal(updateDTO)
		runTest(app, buildPut("/users/"+updatedUser.UserName, jsonRequest)).assertStatus(t, http.StatusNotFound)
	})
}

func buildUpdateUserRequest(user *domain.User) *dto.UpdateUserRequest {
	return &dto.UpdateUserRequest{
		Email:          user.Email,
		Role:           user.Role.String(),
		ConfirmedEmail: user.ConfirmedEmail,
		Phone:          user.Phone,
		Roles:          user.Roles,
	}
}

func buildUsers() []*domain.User {
	return []*domain.User{
		{
			UserName:       "user_id_1",
			Email:          "user1@example.com",
			Role:           role.ADMIN,
			ConfirmedEmail: true,
			Phone:          "12345678",
			Roles:          []string{"viewMember", "viewReports"},
		},
		{
			UserName:       "user_id_2",
			Email:          "user2@example.com",
			Role:           role.USER,
			ConfirmedEmail: false,
			Phone:          "12345678",
			Roles:          []string{"viewMember", "viewReports"},
		},
	}
}
