package user

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/google/uuid"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type Repository interface {
	FindUser(ctx context.Context, username string) (*domain.User, error)
	SaveUser(ctx context.Context, user *domain.User) error
	SearchUser(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.User, error)
}

type dynamoRepository struct {
	*wrapper.DynamoDBWrapper
}

func NewRepository(api wrapper.DynamoDBAPI, userTable string) Repository {
	return &dynamoRepository{
		wrapper.NewDynamoDBWrapper(api, userTable),
	}
}

func (repo dynamoRepository) FindUser(ctx context.Context, username string) (*domain.User, error) {
	resp, err := repo.QueryDynamoDB(ctx, WithUserName(username))
	if err != nil {
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			record := &dto.UserItem{}
			attributevalue.UnmarshalMap(item, record)
			return record.ToUser(), nil
		}
	}
	return nil, nil
}

func (repo dynamoRepository) SaveUser(ctx context.Context, user *domain.User) error {
	user.ID = uuid.NewString()
	return repo.SaveItem(dto.NewUserItem(user))
}

func (repo dynamoRepository) SearchUser(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.User, error) {
	var users = make([]*domain.User, 0)
	resp, err := repo.QueryDynamoDB(ctx, specification)
	if err != nil {
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			record := &dto.UserItem{}
			attributevalue.UnmarshalMap(item, record)
			users = append(users, record.ToUser())
		}
	}
	return users, nil
}
