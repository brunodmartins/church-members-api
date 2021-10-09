package user

import (
	"context"
	"testing"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	mock_wrapper "github.com/BrunoDM2943/church-members-api/platform/aws/wrapper/mock"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	userTable = "User-table"
)

func TestDynamoRepository_FindUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := NewRepository(dynamoMock, userTable)
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(&dynamodb.ScanOutput{Items: buildItems()}, nil)
		user, err := repo.FindUser(userName)
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userName, user.UserName)
		assert.NotNil(t, password, user.Password)
	})
	t.Run("Empty", func(t *testing.T) {
		dynamoMock.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(&dynamodb.ScanOutput{Items: []map[string]types.AttributeValue{}}, nil)
		user, err := repo.FindUser(userName)
		assert.Nil(t, err)
		assert.Nil(t, user)
	})
	t.Run("Error", func(t *testing.T) {
		dynamoMock.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(&dynamodb.ScanOutput{Items: []map[string]types.AttributeValue{}}, genericError)
		user, err := repo.FindUser(userName)
		assert.NotNil(t, err)
		assert.Nil(t, user)
	})
}

func TestDynamoRepository_SaveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := NewRepository(dynamoMock, userTable)
	user := buildUser("", "123")
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			assert.Equal(t, userTable, *params.TableName)
			assert.NotNil(t, params.Item)
			return nil, nil
		})
		err := repo.SaveUser(user)
		assert.Nil(t, err)
		assert.NotEmpty(t, user.ID)
	})
	t.Run("Fail", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, genericError)
		err := repo.SaveUser(user)
		assert.NotNil(t, err)
	})
}

func buildItems() []map[string]types.AttributeValue {
	var items []map[string]types.AttributeValue
	id := domain.NewID()
	items = append(items, buildItem(id))
	return items
}

func buildItem(id string) map[string]types.AttributeValue {
	item, _ := attributevalue.MarshalMap(buildUser(id, ""))
	return item
}
