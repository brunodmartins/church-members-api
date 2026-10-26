package student

import (
	"context"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/google/uuid"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type Repository interface {
	Insert(ctx context.Context, student *domain.Student) error
}

type dynamoRepository struct {
	wrapper *wrapper.DynamoDBWrapper
}

func NewRepository(api wrapper.DynamoDBAPI, table string) Repository {
	return &dynamoRepository{wrapper: wrapper.NewDynamoDBWrapper(api, table)}
}

func (repo dynamoRepository) Insert(ctx context.Context, student *domain.Student) error {
	student.ID = uuid.NewString()
	return repo.wrapper.SaveItem(ctx, dto.NewStudentItem(student))
}
