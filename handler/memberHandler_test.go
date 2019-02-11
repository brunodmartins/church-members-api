package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/member"
	"github.com/gin-gonic/gin"
)

func TestListMembers(t *testing.T) {
	r := gin.Default()
	repo := member.NewMemberInMemoryRepository()
	service := member.NewMemberService(repo)
	memberHandler := NewMemberHandler(service)
	memberHandler.SetUpRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fail()
	}
}

func TestGetMemberBadRequest(t *testing.T) {
	r := gin.Default()
	repo := member.NewMemberInMemoryRepository()
	service := member.NewMemberService(repo)
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
	repo := member.NewMemberInMemoryRepository()
	service := member.NewMemberService(repo)
	memberHandler := NewMemberHandler(service)

	id := entity.NewID()
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
	repo := member.NewMemberInMemoryRepository()
	service := member.NewMemberService(repo)
	memberHandler := NewMemberHandler(service)

	member := &entity.Membro{}
	_, _ = repo.Insert(member)

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
	repo := member.NewMemberInMemoryRepository()
	service := member.NewMemberService(repo)
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
	repo := member.NewMemberInMemoryRepository()
	service := member.NewMemberService(repo)
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
	req, _ := http.NewRequest("POST", "/members", strings.NewReader(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fail()
	}
}

func TestSearch2Results(t *testing.T) {
	r := gin.Default()
	repo := member.NewMemberInMemoryRepository()
	service := member.NewMemberService(repo)
	memberHandler := NewMemberHandler(service)

	_, _ = repo.Insert(&entity.Membro{
		Pessoa: entity.Pessoa{
			Nome: "Bruno",
			Sobrenome: "Damasceno Martins",
		},
	})

	_, _ = repo.Insert(&entity.Membro{
		Pessoa: entity.Pessoa{
			Nome: "Teste",
			Sobrenome: "Brutal",
		},
	})



	memberHandler.SetUpRoutes(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members?q=Bru", nil)
	r.ServeHTTP(w,req)
	result := make([]*entity.Membro, 0)
	if w.Code != http.StatusOK {
		t.Fail()
	}
	_ = json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 2 {
		t.Fail()
	}


}