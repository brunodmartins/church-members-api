package repository

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
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
			record := &dto.MemberItem{}
			attributevalue.UnmarshalMap(item, record)
			members = append(members, record.ToMember())
		}
	}
	return members, nil
}

func (repo dynamoRepository) FindByID(id string) (*model.Member, error) {
	output, err := repo.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value : id,
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
	record := &dto.MemberItem{}
	attributevalue.UnmarshalMap(output.Item, record)
	return record.ToMember(), nil
}

func (repo dynamoRepository) Insert(member *model.Member) (string, error) {
	id := uuid.NewString()
	av, err := attributevalue.MarshalMap(dto.NewMemberItem(member))

	delete(av, "id")
	av["id"] = &types.AttributeValueMemberS{Value: id}


	_, err = repo.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("member"),
	})
	return id, err
}

func (repo dynamoRepository) UpdateStatus(id string, status bool) error {
	_, err := repo.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
		TableName:                 aws.String("member"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":active": &types.AttributeValueMemberBOOL{
				Value: status,
			},
		},
		UpdateExpression:          aws.String("set active = :active"),
	})
	return err
}

func (repo dynamoRepository) GenerateStatusHistory(id string, status bool, reason string, date time.Time) error {
	_, err := repo.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: uuid.New().String()},
			"member_id": &types.AttributeValueMemberS{Value: id},
			"reason": &types.AttributeValueMemberS{Value: reason},
			"status": &types.AttributeValueMemberBOOL{Value: status},
			"date": &types.AttributeValueMemberS{Value: date.String()},
		},
		TableName: aws.String("member_history"),
	})
	return err
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
			record := &dto.MemberItem{}
			attributevalue.UnmarshalMap(item, record)
			members = append(members, record.ToMember())
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
			record := &dto.MemberItem{}
			attributevalue.UnmarshalMap(item, record)
			members = append(members, record.ToMember())
		}
	}
	return members, nil
}