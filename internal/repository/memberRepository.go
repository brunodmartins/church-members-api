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
	FindByID(id string) (*model.Member, error)
	Insert(member *model.Member) (string, error)
	UpdateStatus(ID string, status bool) error
	GenerateStatusHistory(id string, status bool, reason string, date time.Time) error
	FindMembersActive() ([]*model.Member, error)
	FindMembersActiveAndMarried() ([]*model.Member, error)
}

var (
	MemberNotFound = errors.New("Member not found")
)
