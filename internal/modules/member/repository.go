package member

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
	"time"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type Repository interface {
	FindAll(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.Member, error)
	FindByID(ctx context.Context, id string) (*domain.Member, error)
	Insert(ctx context.Context, member *domain.Member) error
	RetireMembership(ctx context.Context, member *domain.Member) error
}

type dynamoRepository struct {
	api         wrapper.DynamoDBAPI
	memberTable string
	wrapper     *wrapper.DynamoDBWrapper
}

func NewRepository(api wrapper.DynamoDBAPI, memberTable string) Repository {
	return dynamoRepository{
		api,
		memberTable,
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

func (repo dynamoRepository) RetireMembership(ctx context.Context, member *domain.Member) error {
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
			":membershipEndDate": &types.AttributeValueMemberS{
				Value: member.MembershipEndDate.Format(time.RFC3339),
			},
			":membershipEndReason": &types.AttributeValueMemberS{
				Value: member.MembershipEndReason,
			},
		},
		ReturnValues:     "UPDATED_NEW",
		UpdateExpression: aws.String("set active = :active, membershipEndDate = :membershipEndDate, membershipEndReason = :membershipEndReason"),
	})
	return err
}
