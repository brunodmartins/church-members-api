package church

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
		result, err := repo.GetByID(nil, id)
		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, table, key, nil, genericError)
		_, err := repo.GetByID(nil, id)
		assert.NotNil(t, err)
	})
	t.Run("Not found", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, table, key, nil, nil)
		_, err := repo.GetByID(nil, id)
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
		result, err := repo.List(nil)
		assert.Nil(t, err)
		assert.Len(t, result, 5)
	})
	t.Run("Empty", func(t *testing.T) {
		wrapper.MockScan(dynamoMock, buildItems(0), nil)
		result, err := repo.List(nil)
		assert.Nil(t, err)
		assert.Len(t, result, 0)
	})
	t.Run("Error", func(t *testing.T) {
		wrapper.MockScan(dynamoMock, buildItems(0), genericError)
		result, err := repo.List(nil)
		assert.NotNil(t, err)
		assert.Len(t, result, 0)
	})
}

func TestDynamoRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	const table = "church"
	ctx := context.TODO()
	repo := NewRepository(dynamoMock, table)

	t.Run("Success", func(t *testing.T) {
		church := buildChurch(uuid.NewString())
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), gomock.Any()).DoAndReturn(
			func(_ context.Context, input *dynamodb.UpdateItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
				assert.Equal(t, table, *input.TableName)
				assert.Equal(t, church.ID, input.Key["id"].(*types.AttributeValueMemberS).Value)
				assert.Equal(t, "attribute_exists(id)", *input.ConditionExpression)
				assert.Equal(t, church.Name, input.ExpressionAttributeValues[":church_name"].(*types.AttributeValueMemberS).Value)
				assert.Equal(t, church.Language, input.ExpressionAttributeValues[":language"].(*types.AttributeValueMemberS).Value)
				assert.Equal(t, church.Email, input.ExpressionAttributeValues[":email"].(*types.AttributeValueMemberS).Value)
				assert.Equal(t, church.Logo, input.ExpressionAttributeValues[":logo"].(*types.AttributeValueMemberS).Value)
				assert.NotContains(t, *input.UpdateExpression, "abbreviation")
				return nil, nil
			},
		)
		assert.NoError(t, repo.Update(ctx, church))
	})

	t.Run("Fail", func(t *testing.T) {
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), gomock.Any()).Return(nil, genericError)
		assert.Error(t, repo.Update(ctx, buildChurch(uuid.NewString())))
	})

	t.Run("Not found", func(t *testing.T) {
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), gomock.Any()).Return(nil, &types.ConditionalCheckFailedException{})
		err := repo.Update(ctx, buildChurch(uuid.NewString()))
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
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
