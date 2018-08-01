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
