package security

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/BrunoDM2943/church-members-api/platform/security/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//go:generate mockgen -source=./auth_repository.go -destination=./mock/auth_repository_mock.go
type UserRepository interface {
	FindUser(username string) (*domain.User, error)
}

type dynamoRepository struct {
	api       wrapper.DynamoDBAPI
	userTable string
}

func NewUserRepository(api wrapper.DynamoDBAPI, userTable string) *dynamoRepository {
	return &dynamoRepository{api: api, userTable: userTable}
}

func (repo dynamoRepository) FindUser(username string) (*domain.User, error) {
	expr := repo.createExpression(username)
	resp, err := repo.api.Scan(context.TODO(), buildScanInput(repo.userTable, expr))
	if err != nil {
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			record := &domain.User{}
			attributevalue.UnmarshalMap(item, record)
			return record, nil
		}
	}
	return nil, nil
}

func buildScanInput(table string, expr expression.Expression) *dynamodb.ScanInput {
	return &dynamodb.ScanInput{
		TableName:                 aws.String(table),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}
}

func (repo dynamoRepository) createExpression(username string) expression.Expression {
	builderExpression := expression.NewBuilder()
	userExpr := expression.Name("username").Equal(expression.Value(username))
	result, _ := builderExpression.WithFilter(userExpr).Build()
	return result
}
