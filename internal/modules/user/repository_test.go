package user

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	tableUser = "User-table"
)

func TestDynamoRepository_FindUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := NewRepository(dynamoMock, tableUser)
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(1)}, nil)
		user, err := repo.FindUser(buildContext(), userName)
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userName, user.UserName)
		assert.NotNil(t, password, user.Password)
	})
	t.Run("Empty", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil)
		user, err := repo.FindUser(buildContext(), userName)
		assert.Nil(t, err)
		assert.Nil(t, user)
	})
	t.Run("Error", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, genericError)
		user, err := repo.FindUser(buildContext(), userName)
		assert.NotNil(t, err)
		assert.Nil(t, user)
	})
}

func TestDynamoRepository_SaveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := NewRepository(dynamoMock, tableUser)
	user := buildUser("", "123")
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			assert.Equal(t, tableUser, *params.TableName)
			assert.NotNil(t, params.Item)
			return nil, nil
		})
		err := repo.SaveUser(buildContext(), user)
		assert.Nil(t, err)
		assert.NotEmpty(t, user.ID)
	})
	t.Run("Fail", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, genericError)
		err := repo.SaveUser(buildContext(), user)
		assert.NotNil(t, err)
	})
}

func TestDynamoRepository_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := NewRepository(dynamoMock, tableUser)
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(2)}, nil)
		members, err := repo.SearchUser(buildContext(), buildMockSpecification(t))
		assert.Nil(t, err)
		assert.Len(t, members, 2)
	})
	t.Run("Empty", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(0)}, nil)
		members, err := repo.SearchUser(buildContext(), buildMockSpecification(t))
		assert.Nil(t, err)
		assert.Len(t, members, 0)
	})
	t.Run("Error", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(0)}, genericError)
		members, err := repo.SearchUser(buildContext(), buildMockSpecification(t))
		assert.NotNil(t, err)
		assert.Len(t, members, 0)
	})
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
	item, _ := attributevalue.MarshalMap(dto.NewUserItem(buildUser(id, "")))
	return item
}

func buildMockSpecification(t *testing.T) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		assert.NotNil(t, builderExpression)
		return wrapper.ExpressionBuilder{
			Builder: builderExpression,
		}
	}
}

func buildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID: "church_id_test",
		},
	})
}
