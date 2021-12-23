package wrapper

import (
	"context"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"net/http"
)

//go:generate mockgen -source=./dynamodb.go -destination=./mock/dynamodb_mock.go
type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

type DynamoDBWrapper struct {
	api   DynamoDBAPI
	table string
}

func NewDynamoDBWrapper(api DynamoDBAPI, table string) *DynamoDBWrapper {
	return &DynamoDBWrapper{api: api, table: table}
}

type QuerySpecification func(ctx context.Context, builderExpression expression.Builder) (string, expression.Builder)

func (wrapper *DynamoDBWrapper) EmptySpecification(table string) QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) (string, expression.Builder) {
		return table, builderExpression
	}
}

func (wrapper *DynamoDBWrapper) ScanDynamoDB(ctx context.Context, specification QuerySpecification) (*dynamodb.ScanOutput, error) {
	builderExpression := expression.NewBuilder()
	table, builderExpression := specification(ctx, builderExpression)

	expr, _ := builderExpression.Build()
	return wrapper.api.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String(table),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	})
}

func (wrapper *DynamoDBWrapper) SaveItem(item interface{}) error {
	av, _ := attributevalue.MarshalMap(item)
	_, err := wrapper.api.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(wrapper.table),
	})
	return err
}

func (wrapper *DynamoDBWrapper) GetItem(key KeyAttribute, value interface{}) error {
	result, err := wrapper.api.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       key.toKeyAttribute(),
		TableName: aws.String(wrapper.table),
	})
	if err != nil {
		return err
	}
	if result.Item == nil {
		return apierrors.NewApiError("Item not found", http.StatusNotFound)
	}

	return attributevalue.UnmarshalMap(result.Item, value)
}

type Key struct {
	Id    string
	Value string
}

type CompositeKey struct {
	PartitionKey Key
	SortKey      Key
}

func (key CompositeKey) toKeyAttribute() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		key.PartitionKey.Id: &types.AttributeValueMemberS{Value: key.PartitionKey.Value},
		key.SortKey.Id:      &types.AttributeValueMemberS{Value: key.SortKey.Value},
	}
}

type PrimaryKey struct {
	Key
}

func (key PrimaryKey) toKeyAttribute() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		key.Id: &types.AttributeValueMemberS{Value: key.Value},
	}
}

type KeyAttribute interface {
	toKeyAttribute() map[string]types.AttributeValue
}
