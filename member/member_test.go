package member

import (
	"testing"

	"github.com/BrunoDM2943/church-members-api/entity"
)

func TestListAllMembers(t *testing.T) {
	repo := NewMemberInMemoryRepository()
	service := NewMemberService(repo)
	repo.dataSet = append(repo.dataSet, &entity.Membro{})
	membros, _ := service.FindAll()
	if len(membros) == 0 {
		t.Error("No members returned from database")
	}
}

func TestFindMember(t *testing.T) {
	repo := NewMemberInMemoryRepository()
	service := NewMemberService(repo)
	id := entity.NewID()
	membro := &entity.Membro{
		ID: id,
	}
	repo.dataSet = append(repo.dataSet, membro)
	membroFound, _ := service.FindByID(id)
	if membroFound.ID != id {
		t.Error("Member not found")
	}
}
