package student

import (
	"context"
	"net/http"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/modules/participant"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
)

//go:generate mockgen -source=./service.go -destination=./mock/service_mock.go
type Service interface {
	CreateStudent(ctx context.Context, student *domain.Student) (string, error)
}

type studentService struct {
	repository         Repository
	memberService      member.Service
	participantService participant.Service
}

func NewService(repository Repository, memberService member.Service, participantService participant.Service) Service {
	return &studentService{
		repository:         repository,
		memberService:      memberService,
		participantService: participantService,
	}
}

func (service studentService) CreateStudent(ctx context.Context, student *domain.Student) (string, error) {
	if student.Member != nil && student.Participant != nil {
		return "", apierrors.NewApiError("A student can be related to either a member or a participant", http.StatusBadRequest)
	}

	if student.Member != nil {
		member, err := service.memberService.GetMember(ctx, student.Member.ID)
		if err != nil {
			return "", err
		}
		if member.Person == nil {
			return "", apierrors.NewApiError("Member has no person information", http.StatusBadRequest)
		}
		student.Member = member
		student.Name = member.Person.GetFullName()
		student.Gender = member.Person.Gender
	}

	if student.Participant != nil {
		participant, err := service.participantService.GetParticipant(ctx, student.Participant.ID)
		if err != nil {
			return "", err
		}
		student.Participant = participant
		student.Name = participant.Name
		student.Gender = participant.Gender
	}

	if student.Member == nil && student.Participant == nil && (student.Name == "" || student.Gender == "") {
		return "", apierrors.NewApiError("Name and gender are required for a new student", http.StatusBadRequest)
	}

	student.ChurchID = domain.GetChurchID(ctx)
	if err := service.repository.Insert(ctx, student); err != nil {
		return "", err
	}
	return student.ID, nil
}
