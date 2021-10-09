package api

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	mock_user "github.com/BrunoDM2943/church-members-api/internal/modules/user/mock"
	"github.com/golang/mock/gomock"
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
		service.EXPECT().SaveUser(gomock.AssignableToTypeOf(new(domain.User))).Return(nil)
		body := getMock("create_user_valid.json")
		runTest(app, buildPost("/users", body)).assertStatus(t, http.StatusCreated)
	})
	t.Run("Fail - Service error - 500", func(t *testing.T) {
		service.EXPECT().SaveUser(gomock.AssignableToTypeOf(new(domain.User))).Return(genericError)
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
