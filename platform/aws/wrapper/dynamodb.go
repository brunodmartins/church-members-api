package wrapper

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//go:generate mockgen -source=./dynamodb.go -destination=./mock/dynamodb_mock.go
type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

type QuerySpecification func(builderExpression expression.Builder) expression.Builder

func ScanDynamoDB(api DynamoDBAPI, specification QuerySpecification, table string) (*dynamodb.ScanOutput, error) {
	builderExpression := expression.NewBuilder()
	builderExpression = specification(builderExpression)

	expr, _ := builderExpression.Build()
	return api.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 aws.String(table),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	})
}

