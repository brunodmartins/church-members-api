package repository

import (
	"errors"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/entity"
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
	FindAll(filters QueryFilters) ([]*entity.Member, error)
	FindByID(id string) (*entity.Member, error)
	Insert(member *entity.Member) (string, error)
	UpdateStatus(ID string, status bool) error
	GenerateStatusHistory(id string, status bool, reason string, date time.Time) error
	FindMembersActive() ([]*entity.Member, error)
	FindMembersActiveAndMarried() ([]*entity.Member, error)
}

var (
	MemberNotFound = errors.New("member not found")
)
