package repository

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

type dynamoRepository struct {
	client *dynamodb.Client
}


func NewDynamoDBRepository() MemberRepository {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	client := dynamodb.NewFromConfig(cfg)
	return dynamoRepository{
		client,
	}
}

func (repo dynamoRepository) FindAll(filters QueryFilters) ([]*model.Member, error) {
	var members = make([]*model.Member, 0)

	builderExpresion := expression.NewBuilder()
	if filters.GetFilter("person.gender") != nil {
		builderExpresion = builderExpresion.WithFilter(expression.Name("gender").Equal(expression.Value(filters.GetFilter("person.gender").(string))))
	}
	if filters.GetFilter("active") != nil {
		builderExpresion = builderExpresion.WithFilter(expression.Name("active").Equal(expression.Value(filters.GetFilter("active").(bool))))
	}
	if filters.GetFilter("name") != nil {
		builderExpresion = builderExpresion.WithFilter(expression.Name("name").Contains(filters.GetFilter("name").(string)))
	}

	expr, _ := builderExpresion.Build()
	resp, err := repo.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String("member"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	})
	if err != nil{
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			members = append(members, unmarshalItem(item))
		}
	}
	return members, nil
}

func (repo dynamoRepository) FindByID(id model.ID) (*model.Member, error) {
	output, err := repo.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value : id.String(),
			},
		},
		TableName: aws.String("member"),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	if output.Item == nil {
		return nil, MemberNotFound
	}
	return unmarshalItem(output.Item), nil
}

func (repo dynamoRepository) Insert(member *model.Member) (model.ID, error) {
	panic("implement me")
}

func (repo dynamoRepository) UpdateStatus(ID model.ID, status bool) error {
	panic("implement me")
}

func (repo dynamoRepository) GenerateStatusHistory(id model.ID, status bool, reason string, date time.Time) error {
	panic("implement me")
}

func (repo dynamoRepository) FindMembersActive() ([]*model.Member, error) {
	var members = make([]*model.Member, 0)

	builderExpresion := expression.NewBuilder().WithFilter(expression.Name("active").Equal(expression.Value(true)))

	expr, _ := builderExpresion.Build()
	resp, err := repo.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String("member"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	})
	if err != nil{
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			members = append(members, unmarshalItem(item))
		}
	}
	return members, nil
}

func (repo dynamoRepository) FindMembersActiveAndMarried() ([]*model.Member, error) {
	var members = make([]*model.Member, 0)

	builderExpresion := expression.NewBuilder()
	builderExpresion = builderExpresion.WithFilter(expression.Name("active").Equal(expression.Value(true)))
	builderExpresion = builderExpresion.WithFilter(expression.Name("marriageDate").AttributeExists())

	expr, _ := builderExpresion.Build()
	resp, err := repo.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String("member"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	})
	if err != nil{
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			members = append(members, unmarshalItem(item))
		}
	}
	return members, nil
}

func unmarshalItem(item map[string]types.AttributeValue) *model.Member {
	member :=  &model.Member{}
	person := &model.Person{}
	contact :=  &model.Contact{}
	address := &model.Address{}
	religion := &model.Religion{}
	attributevalue.UnmarshalMap(item, person)
	attributevalue.UnmarshalMap(item, contact)
	attributevalue.UnmarshalMap(item, address)
	attributevalue.UnmarshalMap(item, religion)
	attributevalue.UnmarshalMap(item, member)

	person.Address = *address
	person.Contact = *contact
	person.Name = person.GetFullName()
	member.Person = *person
	member.Religion = *religion

	return member
}