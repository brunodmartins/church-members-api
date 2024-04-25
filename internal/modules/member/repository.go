package member

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brunodmartins/church-members-api/internal/constants"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/google/uuid"
	"strings"
	"time"
)

//go:generate mockgen -source=./repository.go -destination=./mock/repository_mock.go
type Repository interface {
	FindAll(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.Member, error)
	FindByID(ctx context.Context, id string) (*domain.Member, error)
	Insert(ctx context.Context, member *domain.Member) error
	RetireMembership(ctx context.Context, member *domain.Member) error
	UpdateContact(ctx context.Context, member *domain.Member) error
	UpdateAddress(ctx context.Context, member *domain.Member) error
	UpdatePerson(ctx context.Context, member *domain.Member) error
}

type dynamoRepository struct {
	api         wrapper.DynamoDBAPI
	memberTable string
	wrapper     *wrapper.DynamoDBWrapper
}

func NewRepository(api wrapper.DynamoDBAPI, memberTable string) Repository {
	return dynamoRepository{
		api,
		memberTable,
		wrapper.NewDynamoDBWrapper(api, memberTable),
	}
}

func (repo dynamoRepository) FindAll(ctx context.Context, specification wrapper.QuerySpecification) ([]*domain.Member, error) {
	var members = make([]*domain.Member, 0)
	resp, err := repo.wrapper.QueryDynamoDB(ctx, specification)
	if err != nil {
		return nil, err
	}
	if len(resp.Items) != 0 {
		for _, item := range resp.Items {
			record := &dto.MemberItem{}
			err = attributevalue.UnmarshalMap(item, record)
			if err != nil {
				return nil, err
			}
			members = append(members, record.ToMember())
		}
	}
	return members, nil
}

func (repo dynamoRepository) FindByID(ctx context.Context, id string) (*domain.Member, error) {
	record := &dto.MemberItem{}
	err := repo.wrapper.GetItem(repo.buildKey(ctx, id), record)
	if err != nil {
		return nil, err
	}
	return record.ToMember(), nil
}

func (repo dynamoRepository) buildKey(ctx context.Context, id string) wrapper.CompositeKey {
	return wrapper.CompositeKey{
		PartitionKey: wrapper.Key{
			Id:    "church_id",
			Value: domain.GetChurchID(ctx),
		},
		SortKey: wrapper.Key{
			Id:    "id",
			Value: id,
		},
	}
}

func (repo dynamoRepository) Insert(ctx context.Context, member *domain.Member) error {
	member.ID = uuid.NewString()
	return repo.wrapper.SaveItem(ctx, dto.NewMemberItem(member))
}

func (repo dynamoRepository) RetireMembership(ctx context.Context, member *domain.Member) error {
	_, err := repo.api.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: member.ID,
			},
			"church_id": &types.AttributeValueMemberS{
				Value: member.ChurchID,
			},
		},
		TableName: aws.String(repo.memberTable),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":active": &types.AttributeValueMemberBOOL{
				Value: member.Active,
			},
			":membershipEndDate": &types.AttributeValueMemberS{
				Value: member.MembershipEndDate.Format(time.RFC3339),
			},
			":membershipEndReason": &types.AttributeValueMemberS{
				Value: member.MembershipEndReason,
			},
		},
		ReturnValues:     "UPDATED_NEW",
		UpdateExpression: aws.String(buildUpdateQuery("active", "membershipEndDate", "membershipEndReason")),
	})
	return err
}

func (repo dynamoRepository) UpdateContact(ctx context.Context, member *domain.Member) error {
	_, err := repo.api.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: member.ID,
			},
			"church_id": &types.AttributeValueMemberS{
				Value: member.ChurchID,
			},
		},
		TableName: aws.String(repo.memberTable),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":phoneArea":     toNumberAttributeValue(member.Person.Contact.PhoneArea),
			":phone":         toNumberAttributeValue(member.Person.Contact.Phone),
			":cellPhoneArea": toNumberAttributeValue(member.Person.Contact.CellPhoneArea),
			":cellPhone":     toNumberAttributeValue(member.Person.Contact.CellPhone),
			":email":         toStringAttributeValue(member.Person.Contact.Email),
		},
		ReturnValues:     "UPDATED_NEW",
		UpdateExpression: aws.String(buildUpdateQuery("phoneArea", "phone", "cellPhoneArea", "cellPhone", "email")),
	})
	return err
}

func (repo dynamoRepository) UpdateAddress(ctx context.Context, member *domain.Member) error {
	updateQuery := buildUpdateQuery("zipCode", "state", "city", "address", "district", "number", "moreInfo")
	_, err := repo.api.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: member.ID,
			},
			"church_id": &types.AttributeValueMemberS{
				Value: member.ChurchID,
			},
		},
		TableName: aws.String(repo.memberTable),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":zipCode":  toStringAttributeValue(member.Person.Address.ZipCode),
			":state":    toStringAttributeValue(member.Person.Address.State),
			":city":     toStringAttributeValue(member.Person.Address.City),
			":address":  toStringAttributeValue(member.Person.Address.Address),
			":district": toStringAttributeValue(member.Person.Address.District),
			":number":   toNumberAttributeValue(member.Person.Address.Number),
			":moreInfo": toStringAttributeValue(member.Person.Address.MoreInfo),
		},
		ReturnValues:     "UPDATED_NEW",
		UpdateExpression: aws.String(updateQuery),
	})
	return err
}

func (repo dynamoRepository) UpdatePerson(ctx context.Context, member *domain.Member) error {
	updateQuery := buildUpdateQuery("#name", "firstName", "lastName", "birthDate", "birthDateShort",
		"marriageDate", "marriageDateShort", "spousesName", "maritalStatus", "childrensQuantity")
	updateQuery = strings.Replace(updateQuery, ":#name", ":name", 1)
	attributes := map[string]types.AttributeValue{
		":name":              toStringAttributeValue(member.Person.GetFullName()),
		":firstName":         toStringAttributeValue(member.Person.FirstName),
		":lastName":          toStringAttributeValue(member.Person.LastName),
		":birthDate":         toStringAttributeValue(member.Person.BirthDate.Format(time.RFC3339)),
		":birthDateShort":    toStringAttributeValue(member.Person.BirthDate.Format(constants.ShortDateFormat)),
		":marriageDate":      &types.AttributeValueMemberNULL{Value: true},
		":marriageDateShort": &types.AttributeValueMemberNULL{Value: true},
		":spousesName":       toStringAttributeValue(member.Person.SpousesName),
		":maritalStatus":     toStringAttributeValue(member.Person.MaritalStatus),
		":childrensQuantity": toNumberAttributeValue(member.Person.ChildrenQuantity),
	}
	if member.Person.MarriageDate != nil {
		attributes[":marriageDate"] = toStringAttributeValue(member.Person.MarriageDate.Format(time.RFC3339))
		attributes[":marriageDateShort"] = toStringAttributeValue(member.Person.MarriageDate.Format(constants.ShortDateFormat))
	}
	_, err := repo.api.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: member.ID,
			},
			"church_id": &types.AttributeValueMemberS{
				Value: member.ChurchID,
			},
		},
		TableName: aws.String(repo.memberTable),
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
		},
		ExpressionAttributeValues: attributes,
		ReturnValues:              "UPDATED_NEW",
		UpdateExpression:          aws.String(updateQuery),
	})
	return err
}

func buildUpdateQuery(fields ...string) string {
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

func toNumberAttributeValue(value int) types.AttributeValue {
	if value == 0 {
		return &types.AttributeValueMemberNULL{Value: true}
	}
	return &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", value)}
}
func toStringAttributeValue(value string) types.AttributeValue {
	if value == "" {
		return &types.AttributeValueMemberNULL{Value: true}
	}
	return &types.AttributeValueMemberS{Value: value}
}
