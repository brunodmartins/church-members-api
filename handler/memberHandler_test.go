package handler

import (
	"net/http"
	"net/http/httptest"
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
	_, _ = repo.InsertMember(member)

	memberHandler.SetUpRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members/"+member.ID.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}
