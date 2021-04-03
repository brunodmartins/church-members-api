package repository

import (
	"errors"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/BrunoDM2943/church-members-api/internal/storage/mongo"
)

//go:generate mockgen -source=./memberRepository.go -destination=./mock/memberRepository_mock.go
type MemberRepository interface {
	FindAll(filters mongo.QueryFilters) ([]*model.Member, error)
	FindByID(id model.ID) (*model.Member, error)
	Insert(member *model.Member) (model.ID, error)
	Search(text string) ([]*model.Member, error)
	FindMonthBirthday(date time.Time) ([]*model.Person, error)
	UpdateStatus(ID model.ID, status bool) error
	GenerateStatusHistory(id model.ID, status bool, reason string, date time.Time) error
	FindMembersActive() ([]*model.Member, error)
	FindMembersActiveAndMarried() ([]*model.Member, error)
}

var (
	MemberNotFound = errors.New("Member not found")
)

