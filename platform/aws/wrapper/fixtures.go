package wrapper

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func MockGetItem(t *testing.T, dynamoMock *mock_wrapper.MockDynamoDBAPI, table string, key KeyAttribute, item map[string]types.AttributeValue, err error) {
	dynamoMock.EXPECT().GetItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
		assert.Equal(t, table, *params.TableName)
		assert.Equal(t, key.toKeyAttribute(), params.Key)
		return &dynamodb.GetItemOutput{Item: item}, err
	})
}

func MockQuery(dynamoMock *mock_wrapper.MockDynamoDBAPI, items []map[string]types.AttributeValue, err error) {
	dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: items}, err)
}

func MockScan(dynamoMock *mock_wrapper.MockDynamoDBAPI, items []map[string]types.AttributeValue, err error) {
	dynamoMock.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(&dynamodb.ScanOutput{Items: items}, err)
}
