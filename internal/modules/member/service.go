package member

import (
	"context"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"net/http"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	SearchMembers(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...Specification) ([]*domain.Member, error)
	GetMember(ctx context.Context, id string) (*domain.Member, error)
	SaveMember(ctx context.Context, member *domain.Member) (string, error)
	RetireMembership(ctx context.Context, id string, reason string, date time.Time) error
	UpdateContact(ctx context.Context, memberID string, contact domain.Contact) error
	UpdateAddress(ctx context.Context, memberID string, address domain.Address) error
	UpdatePerson(ctx context.Context, memberID string, person domain.Person) error
}

type memberService struct {
	repo Repository
}

func NewMemberService(r Repository) Service {
	return &memberService{
		repo: r,
	}
}

func (s *memberService) SearchMembers(ctx context.Context, querySpecification wrapper.QuerySpecification, postSpecification ...Specification) ([]*domain.Member, error) {
	members, err := s.repo.FindAll(ctx, querySpecification)
	if err != nil {
		return nil, err
	}
	if len(postSpecification) != 0 {
		return applySpecifications(members, postSpecification), nil
	}
	return members, nil
}

func (s *memberService) GetMember(ctx context.Context, id string) (*domain.Member, error) {
	if !domain.IsValidID(id) {
		return nil, apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	return s.repo.FindByID(ctx, id)
}

func (s *memberService) SaveMember(ctx context.Context, member *domain.Member) (string, error) {
	member.Active = true
	member.ChurchID = domain.GetChurchID(ctx)
	member.MembershipStartDate = time.Now()
	err := s.repo.Insert(ctx, member)
	return member.ID, err
}

func (s *memberService) RetireMembership(ctx context.Context, id string, reason string, date time.Time) error {
	member, err := s.GetMember(ctx, id)
	if err != nil {
		return err
	}
	member.Active = false
	member.MembershipEndDate = &date
	member.MembershipEndReason = reason
	return s.repo.RetireMembership(ctx, member)
}

func (s *memberService) UpdateContact(ctx context.Context, id string, contact domain.Contact) error {
	member, err := s.GetMember(ctx, id)
	if err != nil {
		return err
	}
	member.Person.Contact = &contact
	return s.repo.UpdateContact(ctx, member)
}

func (s *memberService) UpdateAddress(ctx context.Context, id string, address domain.Address) error {
	member, err := s.GetMember(ctx, id)
	if err != nil {
		return err
	}
	member.Person.Address = &address
	return s.repo.UpdateAddress(ctx, member)
}

func (s *memberService) UpdatePerson(ctx context.Context, id string, person domain.Person) error {
	member, err := s.GetMember(ctx, id)
	if err != nil {
		return err
	}
	member.Person.FirstName = person.FirstName
	member.Person.LastName = person.LastName
	member.Person.BirthDate = person.BirthDate
	member.Person.MarriageDate = person.MarriageDate
	member.Person.SpousesName = person.SpousesName
	member.Person.MaritalStatus = person.MaritalStatus
	member.Person.ChildrenQuantity = person.ChildrenQuantity

	return s.repo.UpdatePerson(ctx, member)
}
