package user

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type Repository interface {
	FindUser(username string) (*domain.User, error)
	SaveUser(user *domain.User) error
	SearchUser(specification wrapper.QuerySpecification) ([]*domain.User, error)
}

type dynamoRepository struct {
	api       wrapper.DynamoDBAPI
	userTable string
}

func NewRepository(api wrapper.DynamoDBAPI, userTable string) Repository {
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
			record := &dto.UserItem{}
			attributevalue.UnmarshalMap(item, record)
			return record.ToUser(), nil
		}
	}
	return nil, nil
}

func (repo dynamoRepository) SaveUser(user *domain.User) error {
	user.ID = uuid.NewString()
	av, _ := attributevalue.MarshalMap(dto.NewUserItem(user))
	_, err := repo.api.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repo.userTable),
	})
	return err
}

func (repo dynamoRepository) SearchUser(specification wrapper.QuerySpecification) ([]*domain.User, error) {
	var users = make([]*domain.User, 0)
	resp, err := wrapper.ScanDynamoDB(repo.api, specification, repo.userTable)
	if err != nil {
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			record := &dto.UserItem{}
			attributevalue.UnmarshalMap(item, record)
			users = append(users, record.ToUser())
		}
	}
	return users, nil
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

