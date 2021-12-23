package member_test

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	mock_wrapper "github.com/BrunoDM2943/church-members-api/platform/aws/wrapper/mock"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

const (
	memberTable        = "member-test"
	memberHistoryTable = "member-history-test"
)

func TestDynamoRepository_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := member.NewRepository(dynamoMock, memberTable, memberHistoryTable)
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(&dynamodb.ScanOutput{Items: buildItems(2)}, nil)
		members, err := repo.FindAll(context.Background(), buildMockSpecification(t))
		assert.Nil(t, err)
		assert.Len(t, members, 2)
	})
	t.Run("Empty", func(t *testing.T) {
		dynamoMock.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(&dynamodb.ScanOutput{Items: buildItems(0)}, nil)
		members, err := repo.FindAll(context.Background(), buildMockSpecification(t))
		assert.Nil(t, err)
		assert.Len(t, members, 0)
	})
	t.Run("Error", func(t *testing.T) {
		dynamoMock.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(&dynamodb.ScanOutput{Items: buildItems(0)}, genericError)
		members, err := repo.FindAll(context.Background(), buildMockSpecification(t))
		assert.NotNil(t, err)
		assert.Len(t, members, 0)
	})
}

func TestDynamoRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := member.NewRepository(dynamoMock, memberTable, memberHistoryTable)
	id := domain.NewID()
	ctx := BuildContext()
	t.Run("Success", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, memberTable, buildKey(ctx, id), buildItem(id), nil)
		member, err := repo.FindByID(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, id, member.ID)
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
	repo := member.NewRepository(dynamoMock, memberTable, memberHistoryTable)
	member := buildMember("")
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			assert.Equal(t, memberTable, *params.TableName)
			assert.NotNil(t, params.Item)
			return nil, nil
		})
		err := repo.Insert(context.Background(), member)
		assert.Nil(t, err)
		assert.NotEmpty(t, member.ID)
	})
	t.Run("Fail", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, genericError)
		err := repo.Insert(context.Background(), member)
		assert.NotNil(t, err)
	})
}

func TestDynamoRepository_UpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := member.NewRepository(dynamoMock, memberTable, memberHistoryTable)
	id := domain.NewID()
	member := buildMember(id)
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().UpdateItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			assert.Equal(t, memberTable, *params.TableName)
			assert.Equal(t, id, params.Key["id"].(*types.AttributeValueMemberS).Value)
			return nil, nil
		})
		err := repo.UpdateStatus(context.Background(), member)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		dynamoMock.EXPECT().UpdateItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			assert.Equal(t, memberTable, *params.TableName)
			assert.Equal(t, id, params.Key["id"].(*types.AttributeValueMemberS).Value)
			return nil, genericError
		})
		err := repo.UpdateStatus(context.Background(), member)
		assert.NotNil(t, err)
	})
}

func TestDynamoRepository_GenerateStatusHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := member.NewRepository(dynamoMock, memberTable, memberHistoryTable)
	id := domain.NewID()
	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			assert.Equal(t, memberHistoryTable, *params.TableName)
			return nil, nil
		})
		err := repo.GenerateStatusHistory(id, false, "", time.Now())
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			assert.Equal(t, memberHistoryTable, *params.TableName)
			return nil, genericError
		})
		err := repo.GenerateStatusHistory(id, false, "", time.Now())
		assert.NotNil(t, err)
	})
}

func buildMockSpecification(t *testing.T) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) (string, expression.Builder) {
		assert.NotNil(t, builderExpression)
		return "", builderExpression
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
