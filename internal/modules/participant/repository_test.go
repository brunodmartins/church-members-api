package participant_test

import (
	"context"
	"net/http"
	"testing"

	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/participant"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/brunodmartins/church-members-api/test/dynamodbhelper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	participantTable = "participants-test"
)

func TestDynamoRepository_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := participant.NewRepository(dynamoMock, participantTable)

	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(2)}, nil)
		res, err := repo.FindAll(context.Background(), buildMockSpecification(t))
		assert.Nil(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("Empty", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(0)}, nil)
		res, err := repo.FindAll(context.Background(), buildMockSpecification(t))
		assert.Nil(t, err)
		assert.Len(t, res, 0)
	})

	t.Run("Error", func(t *testing.T) {
		dynamoMock.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&dynamodb.QueryOutput{Items: buildItems(0)}, genericError)
		res, err := repo.FindAll(context.Background(), buildMockSpecification(t))
		assert.NotNil(t, err)
		assert.Len(t, res, 0)
	})
}

func TestDynamoRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := participant.NewRepository(dynamoMock, participantTable)
	id := domain.NewID()
	ctx := BuildContext()

	t.Run("Success", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, participantTable, buildKey(ctx, id), buildItem(id), nil)
		got, err := repo.FindByID(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, id, got.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, participantTable, buildKey(ctx, id), nil, nil)
		got, err := repo.FindByID(ctx, id)
		assert.Equal(t, http.StatusNotFound, err.(apierrors.Error).StatusCode())
		assert.Nil(t, got)
	})

	t.Run("Error", func(t *testing.T) {
		wrapper.MockGetItem(t, dynamoMock, participantTable, buildKey(ctx, id), nil, genericError)
		got, err := repo.FindByID(ctx, id)
		assert.NotNil(t, err)
		assert.Nil(t, got)
	})
}

func TestDynamoRepository_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := participant.NewRepository(dynamoMock, participantTable)
	p := buildParticipant("")

	t.Run("Success", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			assert.Equal(t, participantTable, *params.TableName)
			assert.NotNil(t, params.Item)
			return &dynamodb.PutItemOutput{}, nil
		})
		err := repo.Insert(context.Background(), p)
		assert.Nil(t, err)
		assert.NotEmpty(t, p.ID)
	})

	t.Run("Fail", func(t *testing.T) {
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, genericError)
		err := repo.Insert(context.Background(), p)
		assert.NotNil(t, err)
	})
}

func TestDynamoRepository_UpdateAndDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := participant.NewRepository(dynamoMock, participantTable)

	t.Run("Update Success", func(t *testing.T) {
		p := buildParticipant("u1")
		p.ChurchID = "c1"
		p.Name = "Carol"
		ctx := BuildContext()

		matcher := dynamodbhelper.UpdateMatcher{
			Table:    participantTable,
			ID:       p.ID,
			ChurchID: p.ChurchID,
			Values: map[string]types.AttributeValue{
				":name":        &types.AttributeValueMemberS{Value: p.Name},
				":gender":      &types.AttributeValueMemberNULL{Value: true},
				":birthDate":   &types.AttributeValueMemberS{Value: p.BirthDate.Format(time.RFC3339)},
				":cellPhone":   &types.AttributeValueMemberS{Value: p.CellPhone},
				":filiation":   &types.AttributeValueMemberS{Value: p.Filiation},
				":observation": &types.AttributeValueMemberS{Value: p.Observation},
			},
		}
		dynamoMock.EXPECT().UpdateItem(gomock.Eq(ctx), matcher).Return(nil, nil)
		err := repo.Update(ctx, p)
		assert.Nil(t, err)
	})

	t.Run("Delete Success", func(t *testing.T) {
		id := "d1"
		ctx := BuildContext()
		dynamoMock.EXPECT().DeleteItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
			assert.Equal(t, participantTable, *input.TableName)
			return &dynamodb.DeleteItemOutput{}, nil
		})
		err := repo.Delete(ctx, id)
		assert.Nil(t, err)
	})
}

func buildMockSpecification(t *testing.T) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		assert.NotNil(t, builderExpression)
		return wrapper.ExpressionBuilder{Builder: builderExpression}
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
	item, _ := attributevalue.MarshalMap(dto.NewParticipantItem(buildParticipant(id)))
	return item
}

func buildKey(ctx context.Context, id string) wrapper.CompositeKey {
	return wrapper.CompositeKey{
		PartitionKey: wrapper.Key{Id: "church_id", Value: domain.GetChurchID(ctx)},
		SortKey:      wrapper.Key{Id: "id", Value: id},
	}
}
