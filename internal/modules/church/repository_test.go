package church

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	mock_wrapper "github.com/BrunoDM2943/church-members-api/platform/aws/wrapper/mock"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDynamoRepository_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	const table = "church"
	id := uuid.NewString()
	repo := NewRepository(dynamoMock, table)
	key := repo.(*dynamoRepository).buildKey(id)
	t.Run("Success", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, table, key, buildItem(id), nil)
		result, err := repo.GetByID(id)
		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, table, key, nil, genericError)
		_, err := repo.GetByID(id)
		assert.NotNil(t, err)
	})
	t.Run("Not found", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, table, key, nil, nil)
		_, err := repo.GetByID(id)
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
	})

}

func TestDynamoRepository_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	const table = "church"
	repo := NewRepository(dynamoMock, table)
	t.Run("Success", func(t *testing.T) {
		wrapper.MockScan(dynamoMock, buildItems(5), nil)
		result, err := repo.List()
		assert.Nil(t, err)
		assert.Len(t, result, 5)
	})
	t.Run("Empty", func(t *testing.T) {
		wrapper.MockScan(dynamoMock, buildItems(0), nil)
		result, err := repo.List()
		assert.Nil(t, err)
		assert.Len(t, result, 0)
	})
	t.Run("Error", func(t *testing.T) {
		wrapper.MockScan(dynamoMock, buildItems(0), genericError)
		result, err := repo.List()
		assert.NotNil(t, err)
		assert.Len(t, result, 0)
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
	item, _ := attributevalue.MarshalMap(buildChurch(id))
	return item
}
