package member

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type Repository interface {
	FindAll(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.Member, error)
	FindByID(ctx context.Context, id string) (*domain.Member, error)
	Insert(ctx context.Context, member *domain.Member) error
	UpdateStatus(ctx context.Context, member *domain.Member) error
	GenerateStatusHistory(id string, status bool, reason string, date time.Time) error
}

type dynamoRepository struct {
	api                wrapper.DynamoDBAPI
	memberTable        string
	memberHistoryTable string
	wrapper            *wrapper.DynamoDBWrapper
}

func NewRepository(api wrapper.DynamoDBAPI, memberTable, memberHistoryTable string) Repository {
	return dynamoRepository{
		api,
		memberTable,
		memberHistoryTable,
		wrapper.NewDynamoDBWrapper(api, memberTable),
	}
}

func (repo dynamoRepository) FindAll(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.Member, error) {
	var members = make([]*domain.Member, 0)
	resp, err := repo.wrapper.QueryDynamoDB(ctx, specification)
	if err != nil {
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			record := &dto.MemberItem{}
			attributevalue.UnmarshalMap(item, record)
			members = append(members, record.ToMember())
		}
	}
	return members, nil
}

func (repo dynamoRepository) FindByID(ctx context.Context, id string) (*domain.Member, error) {
	record := &dto.MemberItem{}
	err := repo.wrapper.GetItem(repo.buildKey(ctx, id), record)
	if err != nil {
		return nil, err
	}
	return record.ToMember(), nil
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

func (repo dynamoRepository) Insert(ctx context.Context, member *domain.Member) error {
	member.ID = uuid.NewString()
	return repo.wrapper.SaveItem(dto.NewMemberItem(member))
}

func (repo dynamoRepository) UpdateStatus(ctx context.Context, member *domain.Member) error {
	_, err := repo.api.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: member.ID,
			},
			"church_id": &types.AttributeValueMemberS{
				Value: member.ChurchID,
			},
		},
		TableName: aws.String(repo.memberTable),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":active": &types.AttributeValueMemberBOOL{
				Value: member.Active,
			},
		},
		ReturnValues:     "UPDATED_NEW",
		UpdateExpression: aws.String("set active = :active"),
	})
	return err
}

func (repo dynamoRepository) GenerateStatusHistory(id string, status bool, reason string, date time.Time) error {
	_, err := repo.api.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"id":        &types.AttributeValueMemberS{Value: uuid.New().String()},
			"member_id": &types.AttributeValueMemberS{Value: id},
			"reason":    &types.AttributeValueMemberS{Value: reason},
			"status":    &types.AttributeValueMemberBOOL{Value: status},
			"date":      &types.AttributeValueMemberS{Value: date.String()},
		},
		TableName: aws.String(repo.memberHistoryTable),
	})
	return err
}
