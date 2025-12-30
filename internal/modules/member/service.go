package member

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"

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
	UpdateBaptism(ctx context.Context, memberID string, religion domain.Religion) error
	GetLastBirthAnniversaries(ctx context.Context) ([]*domain.Member, error)
	GetLastMarriageAnniversaries(ctx context.Context) ([]*domain.Member, error)
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

func (s *memberService) UpdateBaptism(ctx context.Context, id string, religion domain.Religion) error {
	member, err := s.GetMember(ctx, id)
	if err != nil {
		return err
	}
	member.Religion.BaptismPlace = religion.BaptismPlace
	member.Religion.Baptized = religion.Baptized
	member.Religion.CatholicBaptized = religion.CatholicBaptized
	member.Religion.BaptismDate = religion.BaptismDate
	return s.repo.UpdateReligion(ctx, member)
}

func (s *memberService) GetLastBirthAnniversaries(ctx context.Context) ([]*domain.Member, error) {
	birthMembers, err := s.SearchMembers(ctx, LastBirths(getWeekRange()))
	if err != nil {
		return nil, err
	}
	sort.Sort(domain.SortByBirthDay(birthMembers))
	return birthMembers, nil
}

func (s *memberService) GetLastMarriageAnniversaries(ctx context.Context) ([]*domain.Member, error) {
	marriageMembers, err := s.SearchMembers(ctx, LastMarriages(getWeekRange()))
	if err != nil {
		return nil, err
	}
	sort.Sort(domain.SortByMarriageDay(marriageMembers))
	return marriageMembers, nil
}

func getWeekRange() (time.Time, time.Time) {
	now := time.Now()

	// Get the current weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	weekday := int(now.Weekday())

	// Calculate days to subtract to get to Monday
	// If Sunday (0), go back 6 days; otherwise go back (weekday - 1) days
	daysToMonday := weekday - 1
	if weekday == 0 {
		daysToMonday = 6
	}

	// Calculate Monday of the current week
	monday := now.AddDate(0, 0, -daysToMonday)

	// Calculate Sunday (6 days after Monday)
	sunday := monday.AddDate(0, 0, 6)

	// Set times to start of day for Monday and end of day for Sunday
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
	sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 999999999, sunday.Location())

	return monday, sunday
}
