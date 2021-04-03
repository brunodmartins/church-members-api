package repository

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/BrunoDM2943/church-members-api/internal/storage/mongo"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

func (repo dynamoRepository) FindAll(filters mongo.QueryFilters) ([]*model.Member, error) {
	panic("implement me")
}

func (repo dynamoRepository) FindByID(id model.ID) (*model.Member, error) {
	output, err := repo.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value : id.String(),
			},
		},
		TableName: aws.String("member"),
	})
	if err != nil {
		return nil, err
	}
	if output.Item != nil {
		return nil, nil
	}
	return UnmarshalItem(output.Item), nil
}

func (repo dynamoRepository) Insert(member *model.Member) (model.ID, error) {
	panic("implement me")
}

func (repo dynamoRepository) Search(text string) ([]*model.Member, error) {
	panic("implement me")
}

func (repo dynamoRepository) FindMonthBirthday(date time.Time) ([]*model.Person, error) {
	panic("implement me")
}

func (repo dynamoRepository) UpdateStatus(ID model.ID, status bool) error {
	panic("implement me")
}

func (repo dynamoRepository) GenerateStatusHistory(id model.ID, status bool, reason string, date time.Time) error {
	panic("implement me")
}

func (repo dynamoRepository) FindMembersActive() ([]*model.Member, error) {
	panic("implement me")
}

func (repo dynamoRepository) FindMembersActiveAndMarried() ([]*model.Member, error) {
	panic("implement me")
}

func UnmarshalItem(item map[string]types.AttributeValue) *model.Member {
	var member model.Member
	var person model.Person
	var contact model.Contact
	var address model.Address
	var religion model.Religion
	attributevalue.UnmarshalMap(item, person)
	attributevalue.UnmarshalMap(item, contact)
	attributevalue.UnmarshalMap(item, address)
	attributevalue.UnmarshalMap(item, religion)
	attributevalue.UnmarshalMap(item, member)

	person.Address = address
	person.Contact = contact
	member.Person = person
	member.Religion = religion

	return &member
}