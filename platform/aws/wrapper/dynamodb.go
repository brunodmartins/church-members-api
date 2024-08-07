package wrapper

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"net/http"
	"strings"
)

//go:generate mockgen -source=./dynamodb.go -destination=./mock/dynamodb_mock.go
type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

type DynamoDBWrapper struct {
	api   DynamoDBAPI
	table string
}

func NewDynamoDBWrapper(api DynamoDBAPI, table string) *DynamoDBWrapper {
	return &DynamoDBWrapper{api: api, table: table}
}

type QuerySpecification func(ctx context.Context, builderExpression expression.Builder) ExpressionBuilder

type ExpressionBuilder struct {
	Index string
	expression.Builder
}

func (wrapper *DynamoDBWrapper) EmptySpecification() QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) ExpressionBuilder {
		return ExpressionBuilder{
			Builder: builderExpression,
		}
	}
}

func (wrapper *DynamoDBWrapper) QueryDynamoDB(ctx context.Context, specification QuerySpecification) (*dynamodb.QueryOutput, error) {
	builderExpression := specification(ctx, expression.NewBuilder())

	expr, _ := builderExpression.Build()
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(wrapper.table),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	if builderExpression.Index != "" {
		queryInput.IndexName = aws.String(builderExpression.Index)
	}
	return wrapper.api.Query(ctx, queryInput)
}

func (wrapper *DynamoDBWrapper) ScanDynamoDB(ctx context.Context, specification QuerySpecification) (*dynamodb.ScanOutput, error) {
	builderExpression := specification(ctx, expression.NewBuilder())

	expr, _ := builderExpression.Build()
	queryInput := &dynamodb.ScanInput{
		TableName:                 aws.String(wrapper.table),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}
	return wrapper.api.Scan(ctx, queryInput)
}

func (wrapper *DynamoDBWrapper) SaveItem(ctx context.Context, item interface{}) error {
	av, _ := attributevalue.MarshalMap(item)
	_, err := wrapper.api.PutItem(ctx, &dynamodb.PutItemInput{
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

func (wrapper *DynamoDBWrapper) BuildUpdateQuery(fields ...string) string {
	builder := strings.Builder{}
	builder.WriteString("set ")
	for _, f := range fields {
		builder.WriteString(f)
		builder.WriteString(" = :")
		builder.WriteString(f)
		builder.WriteString(", ")
	}
	result := builder.String()
	return result[:len(result)-2]
}

// Key groups the ID and Value of a key
type Key struct {
	Id    string
	Value string
}

// CompositeKey groups a partition and sort key
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

// PrimaryKey wraps Key for scenarios where exist only the Partition key
type PrimaryKey struct {
	Key
}

func (key PrimaryKey) toKeyAttribute() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		key.Id: &types.AttributeValueMemberS{Value: key.Value},
	}
}

// KeyAttribute interface provides a conversion from Key to map[string]types.AttributeValue
type KeyAttribute interface {
	toKeyAttribute() map[string]types.AttributeValue
}
