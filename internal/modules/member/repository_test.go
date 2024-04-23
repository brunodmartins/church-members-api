package member_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brunodmartins/church-members-api/internal/constants"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	memberTable = "member-test"
)

func TestDynamoRepository_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := member.NewRepository(dynamoMock, memberTable)
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(2)}, nil)
		members, err := repo.FindAll(context.Background(), buildMockSpecification(t))
		assert.Nil(t, err)
		assert.Len(t, members, 2)
	})
	t.Run("Empty", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(0)}, nil)
		members, err := repo.FindAll(context.Background(), buildMockSpecification(t))
		assert.Nil(t, err)
		assert.Len(t, members, 0)
	})
	t.Run("Error", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(0)}, genericError)
		members, err := repo.FindAll(context.Background(), buildMockSpecification(t))
		assert.NotNil(t, err)
		assert.Len(t, members, 0)
	})
}

func TestDynamoRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := member.NewRepository(dynamoMock, memberTable)
	id := domain.NewID()
	ctx := BuildContext()
	t.Run("Success", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, memberTable, buildKey(ctx, id), buildItem(id), nil)
		churchMember, err := repo.FindByID(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, id, churchMember.ID)
	})
	t.Run("Not Found", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, memberTable, buildKey(ctx, id), nil, nil)
		memberFound, err := repo.FindByID(ctx, id)
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
		assert.Nil(t, memberFound)
	})
	t.Run("Error", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, memberTable, buildKey(ctx, id), nil, genericError)
		memberFound, err := repo.FindByID(ctx, id)
		assert.NotNil(t, err)
		assert.Nil(t, memberFound)
	})
}

func TestDynamoRepository_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := member.NewRepository(dynamoMock, memberTable)
	churchMember := buildMember("")
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			assert.Equal(t, memberTable, *params.TableName)
			assert.NotNil(t, params.Item)
			return nil, nil
		})
		err := repo.Insert(context.Background(), churchMember)
		assert.Nil(t, err)
		assert.NotEmpty(t, churchMember.ID)
	})
	t.Run("Fail", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, genericError)
		err := repo.Insert(context.Background(), churchMember)
		assert.NotNil(t, err)
	})
}

func TestDynamoRepository_UpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := member.NewRepository(dynamoMock, memberTable)
	ctx := context.TODO()

	t.Run("Success - Changing all fields", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		endDate := time.Now()
		churchMember.MembershipEndDate = &endDate
		churchMember.Active = false
		churchMember.MembershipEndReason = "test reason"
		matcher := UpdateMatcher{
			table:    memberTable,
			memberID: churchMember.ID,
			churchID: churchMember.ChurchID,
			values: map[string]types.AttributeValue{
				":active":              &types.AttributeValueMemberBOOL{Value: false},
				":membershipEndDate":   &types.AttributeValueMemberS{Value: endDate.Format(time.RFC3339)},
				":membershipEndReason": &types.AttributeValueMemberS{Value: "test reason"},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, nil)
		err := repo.RetireMembership(ctx, churchMember)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		endDate := time.Now()
		churchMember.MembershipEndDate = &endDate
		churchMember.Active = false
		churchMember.MembershipEndReason = "test reason"
		matcher := UpdateMatcher{
			table:    memberTable,
			memberID: churchMember.ID,
			churchID: churchMember.ChurchID,
			values: map[string]types.AttributeValue{
				":active":              &types.AttributeValueMemberBOOL{Value: false},
				":membershipEndDate":   &types.AttributeValueMemberS{Value: endDate.Format(time.RFC3339)},
				":membershipEndReason": &types.AttributeValueMemberS{Value: "test reason"},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, genericError)
		err := repo.RetireMembership(ctx, churchMember)
		assert.NotNil(t, err)
	})
}

func TestDynamoRepository_UpdateContact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	ctx := context.TODO()
	repo := member.NewRepository(dynamoMock, memberTable)

	t.Run("Success - Changing all fields", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		matcher := UpdateMatcher{
			table:    memberTable,
			memberID: churchMember.ID,
			churchID: churchMember.ChurchID,
			values: map[string]types.AttributeValue{
				":phoneArea":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.Contact.PhoneArea)},
				":phone":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.Contact.Phone)},
				":cellPhoneArea": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.Contact.CellPhoneArea)},
				":cellPhone":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.Contact.CellPhone)},
				":email":         &types.AttributeValueMemberS{Value: churchMember.Person.Contact.Email},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, nil)
		err := repo.UpdateContact(ctx, churchMember)
		assert.Nil(t, err)
	})
	t.Run("Success - Keeping only cellPhone and email", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		churchMember.Person.Contact.PhoneArea = 0
		churchMember.Person.Contact.Phone = 0
		matcher := UpdateMatcher{
			table:    memberTable,
			memberID: churchMember.ID,
			churchID: churchMember.ChurchID,
			values: map[string]types.AttributeValue{
				":phoneArea":     &types.AttributeValueMemberNULL{Value: true},
				":phone":         &types.AttributeValueMemberNULL{Value: true},
				":cellPhoneArea": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.Contact.CellPhoneArea)},
				":cellPhone":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.Contact.CellPhone)},
				":email":         &types.AttributeValueMemberS{Value: churchMember.Person.Contact.Email},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, nil)
		err := repo.UpdateContact(ctx, churchMember)
		assert.Nil(t, err)
	})
	t.Run("Success - Keeping only email", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		churchMember.Person.Contact.PhoneArea = 0
		churchMember.Person.Contact.Phone = 0
		churchMember.Person.Contact.CellPhoneArea = 0
		churchMember.Person.Contact.CellPhone = 0
		matcher := UpdateMatcher{
			table:    memberTable,
			memberID: churchMember.ID,
			churchID: churchMember.ChurchID,
			values: map[string]types.AttributeValue{
				":phoneArea":     &types.AttributeValueMemberNULL{Value: true},
				":phone":         &types.AttributeValueMemberNULL{Value: true},
				":cellPhoneArea": &types.AttributeValueMemberNULL{Value: true},
				":cellPhone":     &types.AttributeValueMemberNULL{Value: true},
				":email":         &types.AttributeValueMemberS{Value: churchMember.Person.Contact.Email},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, nil)
		err := repo.UpdateContact(ctx, churchMember)
		assert.Nil(t, err)
	})
	t.Run("Success - Keeping only cellphone", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		churchMember.Person.Contact.PhoneArea = 0
		churchMember.Person.Contact.Phone = 0
		churchMember.Person.Contact.Email = ""
		matcher := UpdateMatcher{
			table:    memberTable,
			memberID: churchMember.ID,
			churchID: churchMember.ChurchID,
			values: map[string]types.AttributeValue{
				":phoneArea":     &types.AttributeValueMemberNULL{Value: true},
				":phone":         &types.AttributeValueMemberNULL{Value: true},
				":cellPhoneArea": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.Contact.CellPhoneArea)},
				":cellPhone":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.Contact.CellPhone)},
				":email":         &types.AttributeValueMemberNULL{Value: true},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, nil)
		err := repo.UpdateContact(ctx, churchMember)
		assert.Nil(t, err)
	})
	t.Run("Fail - Error on DynamoDB", func(t *testing.T) {
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), gomock.Any()).Return(nil, genericError)
		assert.NotNil(t, repo.UpdateContact(ctx, buildMember(domain.NewID())))
	})
}

func TestDynamoRepository_UpdateAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	ctx := context.TODO()
	repo := member.NewRepository(dynamoMock, memberTable)

	t.Run("Success - Changing all fields", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		matcher := UpdateMatcher{
			table:    memberTable,
			memberID: churchMember.ID,
			churchID: churchMember.ChurchID,
			values: map[string]types.AttributeValue{
				":zipCode":  &types.AttributeValueMemberNULL{Value: true},
				":state":    &types.AttributeValueMemberS{Value: churchMember.Person.Address.State},
				":city":     &types.AttributeValueMemberS{Value: churchMember.Person.Address.City},
				":address":  &types.AttributeValueMemberS{Value: churchMember.Person.Address.Address},
				":district": &types.AttributeValueMemberS{Value: churchMember.Person.Address.District},
				":number":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.Address.Number)},
				":moreInfo": &types.AttributeValueMemberS{Value: churchMember.Person.Address.MoreInfo},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, nil)
		err := repo.UpdateAddress(ctx, churchMember)
		assert.Nil(t, err)
	})
	t.Run("Fail - Error on DynamoDB", func(t *testing.T) {
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), gomock.Any()).Return(nil, genericError)
		assert.NotNil(t, repo.UpdateAddress(ctx, buildMember(domain.NewID())))
	})
}

func TestDynamoRepository_UpdatePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	ctx := context.TODO()
	repo := member.NewRepository(dynamoMock, memberTable)

	t.Run("Success - Changing all fields", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		matcher := UpdateMatcher{
			table:    memberTable,
			memberID: churchMember.ID,
			churchID: churchMember.ChurchID,
			values: map[string]types.AttributeValue{
				":name":              &types.AttributeValueMemberS{Value: churchMember.Person.GetFullName()},
				":firstName":         &types.AttributeValueMemberS{Value: churchMember.Person.FirstName},
				":lastName":          &types.AttributeValueMemberS{Value: churchMember.Person.LastName},
				":birthDate":         &types.AttributeValueMemberS{Value: churchMember.Person.BirthDate.Format(time.RFC3339)},
				":birthDateShort":    &types.AttributeValueMemberS{Value: churchMember.Person.BirthDate.Format(constants.ShortDateFormat)},
				":marriageDate":      &types.AttributeValueMemberS{Value: churchMember.Person.MarriageDate.Format(time.RFC3339)},
				":marriageDateShort": &types.AttributeValueMemberS{Value: churchMember.Person.MarriageDate.Format(constants.ShortDateFormat)},
				":spousesName":       &types.AttributeValueMemberS{Value: churchMember.Person.SpousesName},
				":maritalStatus":     &types.AttributeValueMemberS{Value: churchMember.Person.MaritalStatus},
				":childrensQuantity": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", churchMember.Person.ChildrenQuantity)},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, nil)
		err := repo.UpdatePerson(ctx, churchMember)
		assert.Nil(t, err)
	})
	t.Run("Success - Changing all fields, but no children, nor married", func(t *testing.T) {
		id := domain.NewID()
		churchMember := buildMember(id)
		churchMember.Person.ChildrenQuantity = 0
		churchMember.Person.MaritalStatus = "SINGLE"
		churchMember.Person.MarriageDate = nil
		churchMember.Person.SpousesName = ""
		matcher := UpdateMatcher{
			table:    memberTable,
			memberID: churchMember.ID,
			churchID: churchMember.ChurchID,
			values: map[string]types.AttributeValue{
				":name":              &types.AttributeValueMemberS{Value: churchMember.Person.GetFullName()},
				":firstName":         &types.AttributeValueMemberS{Value: churchMember.Person.FirstName},
				":lastName":          &types.AttributeValueMemberS{Value: churchMember.Person.LastName},
				":birthDate":         &types.AttributeValueMemberS{Value: churchMember.Person.BirthDate.Format(time.RFC3339)},
				":birthDateShort":    &types.AttributeValueMemberS{Value: churchMember.Person.BirthDate.Format(constants.ShortDateFormat)},
				":marriageDate":      &types.AttributeValueMemberNULL{Value: true},
				":marriageDateShort": &types.AttributeValueMemberNULL{Value: true},
				":spousesName":       &types.AttributeValueMemberNULL{Value: true},
				":maritalStatus":     &types.AttributeValueMemberS{Value: churchMember.Person.MaritalStatus},
				":childrensQuantity": &types.AttributeValueMemberNULL{Value: true},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, nil)
		err := repo.UpdatePerson(ctx, churchMember)
		assert.Nil(t, err)
	})
	t.Run("Fail - Error on DynamoDB", func(t *testing.T) {
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), gomock.Any()).Return(nil, genericError)
		assert.NotNil(t, repo.UpdatePerson(ctx, buildMember(domain.NewID())))
	})
}

type UpdateMatcher struct {
	table    string
	memberID string
	churchID string
	values   map[string]types.AttributeValue
}

func (expected UpdateMatcher) Matches(r any) bool {
	received := r.(*dynamodb.UpdateItemInput)
	if *received.TableName != expected.table {
		return false
	}
	if received.Key["id"].(*types.AttributeValueMemberS).Value != expected.memberID {
		return false
	}
	if received.Key["church_id"].(*types.AttributeValueMemberS).Value != expected.churchID {
		return false
	}
	if received.UpdateExpression == nil {
		return false
	}
	updateQuery := *received.UpdateExpression
	if !reflect.DeepEqual(received.ExpressionAttributeValues, expected.values) {
		return false
	}
	for key := range received.ExpressionAttributeValues {
		if !strings.Contains(updateQuery, key) {
			log.Errorf("Key %s not found in update query", key)
			return false
		}
	}

	return true
}

func (expected UpdateMatcher) String() string {
	return fmt.Sprintf("Expected ID: {%s}, ChurchID: {%s}, Table: {%s}, Values:{%v}", expected.memberID, expected.churchID, expected.table, expected.values)
}

func buildMockSpecification(t *testing.T) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		assert.NotNil(t, builderExpression)
		return wrapper.ExpressionBuilder{
			Builder: builderExpression,
		}
	}
}

func buildItems(length int) []map[string]types.AttributeValue {
	var items []map[string]types.AttributeValue
	for i := 0; i < length; i++ {
		id := domain.NewID()
		items = append(items, buildItem(id))
	}
	return items
}

func buildItem(id string) map[string]types.AttributeValue {
	item, _ := attributevalue.MarshalMap(buildMember(id))
	return item
}

func buildKey(ctx context.Context, id string) wrapper.CompositeKey {
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
