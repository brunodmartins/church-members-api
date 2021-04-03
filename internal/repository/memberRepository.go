package repository

import (
	"errors"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
)

type QueryFilters map[string]interface{}

func (qf QueryFilters) AddFilter(key string, value interface{}) {
	qf[key] = value
}

func (qf QueryFilters) GetFilter(key string) interface{} {
	return qf[key]
}

//go:generate mockgen -source=./memberRepository.go -destination=./mock/memberRepository_mock.go
type MemberRepository interface {
	FindAll(filters QueryFilters) ([]*model.Member, error)
	FindByID(id model.ID) (*model.Member, error)
	Insert(member *model.Member) (model.ID, error)
	UpdateStatus(ID model.ID, status bool) error
	GenerateStatusHistory(id model.ID, status bool, reason string, date time.Time) error
	FindMembersActive() ([]*model.Member, error)
	FindMembersActiveAndMarried() ([]*model.Member, error)
}

var (
	MemberNotFound = errors.New("Member not found")
)

