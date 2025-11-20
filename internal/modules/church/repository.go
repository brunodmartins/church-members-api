package church

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type Repository interface {
	GetByID(ctx context.Context, ID string) (*domain.Church, error)
	List(ctx context.Context) ([]*domain.Church, error)
}

type dynamoRepository struct {
	*wrapper.DynamoDBWrapper
}

func NewRepository(api wrapper.DynamoDBAPI, table string) Repository {
	return &dynamoRepository{
		wrapper.NewDynamoDBWrapper(api, table),
	}
}

func (d dynamoRepository) GetByID(ctx context.Context, id string) (*domain.Church, error) {
	result := &domain.Church{}
	return result, d.GetItem(d.buildKey(id), result)
}

func (d dynamoRepository) buildKey(id string) wrapper.PrimaryKey {
	return wrapper.PrimaryKey{
		Key: wrapper.Key{
			Id:    "id",
			Value: id,
		}}
}

func (d dynamoRepository) List(ctx context.Context) ([]*domain.Church, error) {
	var result = make([]*domain.Church, 0)
	resp, err := d.ScanDynamoDB(ctx, d.EmptySpecification())
	if err != nil {
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			record := &domain.Church{}
			attributevalue.UnmarshalMap(item, record)
			result = append(result, record)
		}
	}
	return result, nil
}
