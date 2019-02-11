package member

import "github.com/BrunoDM2943/church-members-api/entity"

type MemberService struct {
	repo Repository
}

func NewMemberService(r Repository) *MemberService {
	return &MemberService{
		repo: r,
	}
}

func (s *MemberService) FindAll() ([]*entity.Membro, error) {
	return s.repo.FindAll()
}

func (s *MemberService) FindByID(id entity.ID) (*entity.Membro, error) {
	return s.repo.FindByID(id)
}

func (s *MemberService) Insert(membro *entity.Membro) (entity.ID, error) {
	return s.repo.Insert(membro)
}

func (s *MemberService) Search(text string) ([]*entity.Membro, error) {
	return s.repo.Search(text)
}