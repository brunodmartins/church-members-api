package user

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/google/uuid"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type Repository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindUser(ctx context.Context, username string) (*domain.User, error)
	SaveUser(ctx context.Context, user *domain.User) error
	SearchUser(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
}

type dynamoRepository struct {
	api wrapper.DynamoDBAPI
	*wrapper.DynamoDBWrapper
	table string
}

func NewRepository(api wrapper.DynamoDBAPI, userTable string) Repository {
	return &dynamoRepository{
		api,
		wrapper.NewDynamoDBWrapper(api, userTable),
		userTable,
	}
}

func (repo dynamoRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	record := &dto.UserItem{}
	err := repo.GetItem(repo.buildKey(ctx, id), record)
	if err != nil {
		return nil, err
	}
	return record.ToUser(), nil
}

func (repo dynamoRepository) FindUser(ctx context.Context, username string) (*domain.User, error) {
	resp, err := repo.QueryDynamoDB(ctx, WithUserName(username))
	if err != nil {
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			record := &dto.UserItem{}
			err = attributevalue.UnmarshalMap(item, record)
			if err != nil {
				return nil, err
			}
			return record.ToUser(), nil
		}
	}
	return nil, nil
}

func (repo dynamoRepository) SaveUser(ctx context.Context, user *domain.User) error {
	user.ID = uuid.NewString()
	return repo.SaveItem(ctx, dto.NewUserItem(user))
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
			err = attributevalue.UnmarshalMap(item, record)
			if err != nil {
				return nil, err
			}
			users = append(users, record.ToUser())
		}
	}
	return users, nil
}

func (repo dynamoRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := repo.api.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: user.ID,
			},
			"church_id": &types.AttributeValueMemberS{
				Value: user.ChurchID,
			},
		},
		TableName: aws.String(repo.table),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":confirmed_email": &types.AttributeValueMemberBOOL{Value: user.ConfirmedEmail},
		},
		ReturnValues:     "UPDATED_NEW",
		UpdateExpression: aws.String(repo.BuildUpdateQuery("confirmed_email")),
	})
	return err
}

func (repo dynamoRepository) buildKey(ctx context.Context, id string) wrapper.CompositeKey {
	return wrapper.CompositeKey{
		PartitionKey: wrapper.Key{
			Id:    "church_id",
			Value: domain.GetChurchID(ctx),
		},
		SortKey: wrapper.Key{
			Id:    "id",
			Value: id,
		},
	}
}
