package participant

import (
	"context"
	"time"

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
	Insert(ctx context.Context, p *domain.Participant) error
	GetByID(ctx context.Context, id string) (*domain.Participant, error)
	Update(ctx context.Context, p *domain.Participant) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.Participant, error)
}

type dynamoRepository struct {
	api     wrapper.DynamoDBAPI
	table   string
	wrapper *wrapper.DynamoDBWrapper
}

func NewRepository(api wrapper.DynamoDBAPI, table string) Repository {
	return &dynamoRepository{
		api:     api,
		table:   table,
		wrapper: wrapper.NewDynamoDBWrapper(api, table),
	}
}

func (repo dynamoRepository) Insert(ctx context.Context, participant *domain.Participant) error {
	participant.ID = uuid.NewString()
	return repo.wrapper.SaveItem(ctx, dto.NewParticipantItem(participant))
}

func (repo dynamoRepository) GetByID(ctx context.Context, id string) (*domain.Participant, error) {
	record := &dto.ParticipantItem{}
	err := repo.wrapper.GetItem(repo.buildKey(ctx, id), record)
	if err != nil {
		return nil, err
	}
	return record.ToParticipant(), nil
}

func (repo dynamoRepository) Update(ctx context.Context, participant *domain.Participant) error {
	_, err := repo.api.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: participant.ID,
			},
			"church_id": &types.AttributeValueMemberS{
				Value: participant.ChurchID,
			},
		},
		TableName: aws.String(repo.table),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":   toStringAttributeValue(participant.Name),
			":gender": toStringAttributeValue(participant.Gender),
			":birthDate": &types.AttributeValueMemberS{
				Value: participant.BirthDate.Format(time.RFC3339),
			},
			":cellPhone":   toStringAttributeValue(participant.CellPhone),
			":filiation":   toStringAttributeValue(participant.Filiation),
			":observation": toStringAttributeValue(participant.Observation),
		},
		ReturnValues:     "UPDATED_NEW",
		UpdateExpression: aws.String(repo.wrapper.BuildUpdateQuery("name", "gender", "birthDate", "cellPhone", "filiation", "observation")),
	})
	return err
}

func (repo dynamoRepository) Delete(ctx context.Context, id string) error {
	_, err := repo.api.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(repo.table),
		Key: map[string]types.AttributeValue{
			"church_id": &types.AttributeValueMemberS{Value: domain.GetChurchID(ctx)},
			"id":        &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}

func (repo dynamoRepository) FindAll(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.Participant, error) {
	var result = make([]*domain.Participant, 0)
	resp, err := repo.wrapper.QueryDynamoDB(ctx, specification)
	if err != nil {
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			record := &dto.ParticipantItem{}
			err = attributevalue.UnmarshalMap(item, record)
			if err != nil {
				return nil, err
			}
			result = append(result, record.ToParticipant())
		}
	}
	return result, nil
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

func toStringAttributeValue(value string) types.AttributeValue {
	if value == "" {
		return &types.AttributeValueMemberNULL{Value: true}
	}
	return &types.AttributeValueMemberS{Value: value}
}
