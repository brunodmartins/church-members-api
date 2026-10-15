package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	mock_church "github.com/brunodmartins/church-members-api/internal/modules/church/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"go.uber.org/mock/gomock"
)

func TestGetChurch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_church.NewMockService(ctrl)
	handler := NewChurchHandler(service)
	handler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		id := domain.NewID()
		church := buildChurch(id)
		service.EXPECT().GetChurch(gomock.Any(), id).Return(church, nil)
		expected := dto.GetChurchResponse{
			ID:           church.ID,
			Name:         church.Name,
			Abbreviation: church.Abbreviation,
			Language:     church.Language,
			Email:        church.Email,
			Logo:         church.Logo,
		}
		runTest(app, buildGet("/churches/"+id)).assertJSON(t, http.StatusOK, expected)
	})
	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetChurch(gomock.Any(), id).Return(nil, apierrors.NewApiError("Member not found", http.StatusNotFound))
		runTest(app, buildGet("/churches/"+id)).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetChurch(gomock.Any(), id).Return(nil, genericError)
		runTest(app, buildGet("/churches/"+id)).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestGetStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_church.NewMockService(ctrl)
	handler := NewChurchHandler(service)
	handler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		id := domain.NewID()
		stats := buildStatistics()
		service.EXPECT().GetStatistics(gomock.Any(), id).Return(stats, nil)
		expected := dto.ChurchStatisticsResponse{
			TotalMembers:                 stats.TotalMembers,
			AgeDistribution:              stats.AgeDistribution,
			TotalMembersByGender:         stats.TotalMembersByGender,
			TotalMembersByClassification: stats.TotalMembersByClassification,
		}
		runTest(app, buildGet("/churches/"+id+"/statistics")).assertJSON(t, http.StatusOK, expected)
	})
	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetStatistics(gomock.Any(), id).Return(nil, apierrors.NewApiError("Church not found", http.StatusNotFound))
		runTest(app, buildGet("/churches/"+id+"/statistics")).assertStatus(t, http.StatusNotFound)
	})
	t.Run("Fail - 500", func(t *testing.T) {
		id := domain.NewID()
		service.EXPECT().GetStatistics(gomock.Any(), id).Return(nil, genericError)
		runTest(app, buildGet("/churches/"+id+"/statistics")).assertStatus(t, http.StatusInternalServerError)
	})
}

func TestUpdateChurch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := newApp()

	service := mock_church.NewMockService(ctrl)
	handler := NewChurchHandler(service)
	handler.SetUpRoutes(app)

	t.Run("Success - 200", func(t *testing.T) {
		id := domain.NewID()
		request := &dto.UpdateChurchRequest{
			Name:     "updated church",
			Language: "en-us",
			Email:    "church@example.com",
		}
		expected := request.ToChurch()
		expected.ID = id
		service.EXPECT().UpdateChurch(gomock.Any(), gomock.Eq(expected)).Return(nil)
		jsonRequest, _ := json.Marshal(request)
		runTest(app, buildPut("/churches/"+id, jsonRequest)).assertJSON(t, http.StatusOK, dto.MessageResponse{Message: "Church updated successfully"})
	})

	t.Run("Fail - invalid ID - 400", func(t *testing.T) {
		jsonRequest, _ := json.Marshal(&dto.UpdateChurchRequest{Name: "updated church", Language: "en-us"})
		runTest(app, buildPut("/churches/invalid-id", jsonRequest)).assertStatus(t, http.StatusBadRequest)
	})

	t.Run("Fail - invalid body - 400", func(t *testing.T) {
		id := domain.NewID()
		jsonRequest, _ := json.Marshal(&dto.UpdateChurchRequest{Language: "en-us"})
		runTest(app, buildPut("/churches/"+id, jsonRequest)).assertStatus(t, http.StatusBadRequest)
	})

	t.Run("Fail - malformed JSON - 400", func(t *testing.T) {
		id := domain.NewID()
		runTest(app, buildPut("/churches/"+id, badJson)).assertStatus(t, http.StatusBadRequest)
	})

	t.Run("Fail - 404", func(t *testing.T) {
		id := domain.NewID()
		request := &dto.UpdateChurchRequest{Name: "updated church", Language: "en-us"}
		expected := request.ToChurch()
		expected.ID = id
		service.EXPECT().UpdateChurch(gomock.Any(), gomock.Eq(expected)).Return(apierrors.NewApiError("Church not found", http.StatusNotFound))
		jsonRequest, _ := json.Marshal(request)
		runTest(app, buildPut("/churches/"+id, jsonRequest)).assertStatus(t, http.StatusNotFound)
	})

	t.Run("Fail - 403", func(t *testing.T) {
		id := domain.NewID()
		request := &dto.UpdateChurchRequest{Name: "updated church", Language: "en-us"}
		expected := request.ToChurch()
		expected.ID = id
		service.EXPECT().UpdateChurch(gomock.Any(), gomock.Eq(expected)).Return(apierrors.NewApiError("User does not have required role", http.StatusForbidden))
		jsonRequest, _ := json.Marshal(request)
		runTest(app, buildPut("/churches/"+id, jsonRequest)).assertStatus(t, http.StatusForbidden)
	})
}

func TestUpdateChurchLogo(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mock_church.NewMockService(ctrl)
	app := newApp()
	NewChurchHandler(service).SetUpRoutes(app)

	id := domain.NewID()
	pngData := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	t.Run("Success - 200", func(t *testing.T) {
		service.EXPECT().UpdateChurchLogo(gomock.Any(), id, pngData).Return(nil)
		runTest(app, buildLogoPost("/churches/"+id+"/logo", pngData)).assertJSON(
			t,
			http.StatusOK,
			dto.MessageResponse{Message: "Church logo updated successfully"},
		)
	})

	t.Run("Fail - invalid ID - 400", func(t *testing.T) {
		runTest(app, buildLogoPost("/churches/invalid-id/logo", pngData)).assertStatus(t, http.StatusBadRequest)
	})

	t.Run("Fail - missing logo - 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/churches/"+id+"/logo", nil)
		req.Header.Set("Content-Type", "multipart/form-data")
		runTest(app, req).assertStatus(t, http.StatusBadRequest)
	})

	t.Run("Fail - non-PNG file - 400", func(t *testing.T) {
		runTest(app, buildLogoPost("/churches/"+id+"/logo", []byte("not a PNG"))).assertStatus(t, http.StatusBadRequest)
	})

	t.Run("Fail - file is 5 MB - 400", func(t *testing.T) {
		runTest(app, buildLogoPost("/churches/"+id+"/logo", bytes.Repeat([]byte{'a'}, maxChurchLogoSize))).assertStatus(t, http.StatusBadRequest)
	})
}

func buildLogoPost(target string, data []byte) *http.Request {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("logo", "logo.png")
	_, _ = part.Write(data)
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, target, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func buildStatistics() *domain.ChurchStatistics {
	return &domain.ChurchStatistics{
		TotalMembers:                 10,
		AgeDistribution:              []int{1, 2, 3},
		TotalMembersByGender:         map[string]int{"M": 5, "F": 5},
		TotalMembersByClassification: map[string]int{"MEMBER": 10},
	}
}
