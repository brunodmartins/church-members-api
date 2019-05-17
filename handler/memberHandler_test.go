package handler

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/member/repository"
	mock_service "github.com/BrunoDM2943/church-members-api/member/service/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	member := &entity.Membro{}
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
         "dddTelefone": 11,
         "telefone": 29435002,
         "dddCelular": 11,
         "celular": 953200587,
         "email": "bdm2943@gmail.com"
      },
      "endereco": {
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
         "dddTelefone": 11,
         "telefone": 29435002,
         "dddCelular": 11,
         "celular": 953200587,
         "email": "bdm2943@gmail.com"
      },
      "endereco": {
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
	service.EXPECT().FindMembers(gomock.Any()).Return([]*entity.Membro{}, nil)
	req, _ := http.NewRequest("POST", "/members/search", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

