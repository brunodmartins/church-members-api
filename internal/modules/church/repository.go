package church

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type Repository interface {
	GetByID(ID string) (*domain.Church, error)
	List() ([]*domain.Church, error)
}

type dynamoRepository struct {
	*wrapper.DynamoDBWrapper
}

func NewRepository(api wrapper.DynamoDBAPI, table string) Repository {
	return &dynamoRepository{
		wrapper.NewDynamoDBWrapper(api, table),
	}
}

func (d dynamoRepository) GetByID(id string) (*domain.Church, error) {
	result := &domain.Church{}
	return result, d.GetItem(id, result)
}

func (d dynamoRepository) List() ([]*domain.Church, error) {
	var result = make([]*domain.Church, 0)
	resp, err := d.ScanDynamoDB(context.Background(), d.EmptySpecification())
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
