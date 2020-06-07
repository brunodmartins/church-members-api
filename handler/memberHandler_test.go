package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BrunoDM2943/church-members-api/member/repository"
	mock_service "github.com/BrunoDM2943/church-members-api/member/service/mock"

	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestGetMemberBadRequest(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockIMemberService(ctrl)
	memberHandler := NewMemberHandler(service)

	memberHandler.SetUpRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members/a", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestGetMemberNotFound(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockIMemberService(ctrl)
	memberHandler := NewMemberHandler(service)

	id := entity.NewID()

	service.EXPECT().FindMembersByID(id).Return(nil, repository.MemberNotFound)
	memberHandler.SetUpRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fail()
	}
}

func TestGetMemberOK(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockIMemberService(ctrl)
	memberHandler := NewMemberHandler(service)

	member := &entity.Member{}
	member.ID = entity.NewID()
	service.EXPECT().FindMembersByID(member.ID).Return(member, nil).AnyTimes()

	memberHandler.SetUpRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members/"+member.ID.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestPostMemberBadRequest(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockIMemberService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/members", strings.NewReader(""))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestPostMemberSucess(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockIMemberService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	body := `
	{
   "pessoa": {
      "contato": {
         "dddPhone": 11,
         "telefone": 29435002,
         "dddCellPhone": 11,
         "celular": 953200587,
         "email": "bdm2943@gmail.com"
      },
      "address": {
         "cep": "03805090",
         "uf": "SP",
         "cidade": "São Paulo",
         "logradouro": "Rua Dario Costa Mattos",
         "bairro": "Parque Boturussu",
         "numero": 661,
         "complemento": "Casa"
      },
      "escolaridade": {
         "ensinoFundamental": true,
         "ensinoMedio": true,
         "ensinoSuperior": true
      },
      "nome": "Bruno",
      "sobrenome": "Damasceno",
      "dtNascimento": "1995-05-10T00:00:00-03:00",
      "naturalidade": "Brasil",
      "cidadeNascimento": "São Paulo",
      "nomeMae": "Mae",
      "nomePai": "Pai",
      "estadoCivil": "S",
      "sexo":"M",
      "qtdIrmao": 1,
      "qtdFilhos": 1,
      "profissao": "Teste"
   },
   "religiao": {
      "religiaoPais": "Crentes",
      "batizadoCatolica": true,
      "idadeConheceuEvangelho": 10,
      "aceitouJesus": true,
      "dtAceitouJesus": null,
      "batizado": true,
      "dtBatismo": null,
      "localBatismo": "IEPEM",
      "conheceDizimo": true,
      "concordaDizimo": true,
      "dizimista": true
   },
   "frequentaCultoSexta": true,
   "frequentaCultoSabado": true,
   "frequentaEBD": true,
   "frequentaCultoDomingo": true
}`
	service.EXPECT().SaveMember(gomock.Any()).Return(entity.NewID(), nil)
	req, _ := http.NewRequest("POST", "/members", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fail()
	}
}

func TestPostMemberFail(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockIMemberService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	body := `
	{
   "pessoa": {
      "contato": {
         "dddPhone": 11,
         "telefone": 29435002,
         "dddCellPhone": 11,
         "celular": 953200587,
         "email": "bdm2943@gmail.com"
      },
      "address": {
         "cep": "03805090",
         "uf": "SP",
         "cidade": "São Paulo",
         "logradouro": "Rua Dario Costa Mattos",
         "bairro": "Parque Boturussu",
         "numero": 661,
         "complemento": "Casa"
      },
      "escolaridade": {
         "ensinoFundamental": true,
         "ensinoMedio": true,
         "ensinoSuperior": true
      },
      "nome": "Bruno",
      "sobrenome": "Damasceno",
      "dtNascimento": "1995-05-10T00:00:00-03:00",
      "naturalidade": "Brasil",
      "cidadeNascimento": "São Paulo",
      "nomeMae": "Mae",
      "nomePai": "Pai",
      "estadoCivil": "S",
      "sexo":"M",
      "qtdIrmao": 1,
      "qtdFilhos": 1,
      "profissao": "Teste"
   },
   "religiao": {
      "religiaoPais": "Crentes",
      "batizadoCatolica": true,
      "idadeConheceuEvangelho": 10,
      "aceitouJesus": true,
      "dtAceitouJesus": null,
      "batizado": true,
      "dtBatismo": null,
      "localBatismo": "IEPEM",
      "conheceDizimo": true,
      "concordaDizimo": true,
      "dizimista": true
   },
   "frequentaCultoSexta": true,
   "frequentaCultoSabado": true,
   "frequentaEBD": true,
   "frequentaCultoDomingo": true
}`
	service.EXPECT().SaveMember(gomock.Any()).Return(entity.NewID(), errors.New(""))
	req, _ := http.NewRequest("POST", "/members", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fail()
	}
}

func TestPostMemberSearch(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockIMemberService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	body := `
	{
		member(sexo:"M", active:false){
				pessoa{
					nome,
					sobrenome
					sexo
				}
		}
	}`
	service.EXPECT().FindMembers(gomock.Any()).Return([]*entity.Member{}, nil)
	req, _ := http.NewRequest("POST", "/members/search", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestPostMemberSearchError(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockIMemberService(ctrl)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	body := `
	{
		member(sexo:"M", active:false){
				pessoa{
					nome,
					sobreno
					sexo
				}
		}
	}`
	req, _ := http.NewRequest("POST", "/members/search", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fail()
	}
}

func TestPutStatus(t *testing.T) {
	type data struct {
		url        string
		body       string
		statusCode int
	}
	id := entity.NewID()
	urlWithID := "/members/" + id.String() + "/status"
	table := []data{
		{"/members/X/status", "", http.StatusBadRequest},
		{urlWithID, `{"active":false}`, http.StatusBadRequest},
		{urlWithID, `{"reason": "exited"}`, http.StatusBadRequest},
		{urlWithID, `{"active":false, "reason": "exited"}`, http.StatusInternalServerError},
		{urlWithID, `{"active":true, "reason": "Comed back"}`, http.StatusOK}}

	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockIMemberService(ctrl)
	memberHandler := NewMemberHandler(service)
	service.EXPECT().ChangeStatus(id, false, "exited").Return(errors.New("Error"))
	service.EXPECT().ChangeStatus(id, true, "Comed back").Return(nil)

	memberHandler.SetUpRoutes(r)

	for _, test := range table {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", test.url, strings.NewReader(test.body))
		r.ServeHTTP(w, req)
		if w.Code != test.statusCode {
			t.Errorf("Failed for test: %s, %s, %d", test.url, test.body, test.statusCode)
		}
	}
}
