package participant

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDynamoRepository_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dynamoMock := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	const table = "participants"
	repo := NewRepository(dynamoMock, table)

	// Insert
	t.Run("Insert Success", func(t *testing.T) {
		ctx := BuildContext()
		p := buildParticipant(domain.NewID())
		p.Name = "Alice"
		dynamoMock.EXPECT().PutItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, input *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				assert.Equal(t, table, *input.TableName)
				return &dynamodb.PutItemOutput{}, nil
			})
		err := repo.Insert(ctx, p)
		assert.NoError(t, err)
		assert.NotEmpty(t, p.ID)
	})

	// GetByID
	t.Run("GetByID Success", func(t *testing.T) {
		id := uuid.NewString()
		p := buildParticipant(id)
		p.ChurchID = "c1"
		p.Name = "Bob"
		item := dto.NewParticipantItem(p)
		itemMap, _ := attributevalue.MarshalMap(item)

		ctx := BuildContext()
		key := repo.(*dynamoRepository).buildKey(ctx, id)
		wrapper.MockGetItem(t, dynamoMock, table, key, itemMap, nil)

		got, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, "Bob", got.Name)
	})

	t.Run("GetByID NotFound", func(t *testing.T) {
		id := uuid.NewString()
		ctx := BuildContext()
		key := repo.(*dynamoRepository).buildKey(ctx, id)
		wrapper.MockGetItem(t, dynamoMock, table, key, nil, nil)
		_, err := repo.GetByID(ctx, id)
		assert.NotNil(t, err)
		apiErr, ok := err.(apierrors.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	// Update
	t.Run("Update Success", func(t *testing.T) {
		p := buildParticipant("u1")
		p.ChurchID = "c1"
		p.Name = "Carol"
		ctx := BuildContext()
		dynamoMock.EXPECT().UpdateItem(gomock.Any(), gomock.Any(), gomock.Any()).Return(&dynamodb.UpdateItemOutput{}, nil)
		err := repo.Update(ctx, p)
		assert.NoError(t, err)
	})

	// Delete
	t.Run("Delete Success", func(t *testing.T) {
		id := "d1"
		ctx := BuildContext()
		dynamoMock.EXPECT().DeleteItem(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, input *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
				assert.Equal(t, table, *input.TableName)
				return &dynamodb.DeleteItemOutput{}, nil
			})
		err := repo.Delete(ctx, id)
		assert.NoError(t, err)
	})

	// FindAll / Query
	t.Run("FindAll Success", func(t *testing.T) {
		p1 := buildParticipant("q1")
		p1.ChurchID = "c1"
		p1.Name = "Q1"
		p2 := buildParticipant("q2")
		p2.ChurchID = "c1"
		p2.Name = "Q2"
		i1 := dto.NewParticipantItem(p1)
		i2 := dto.NewParticipantItem(p2)
		m1, _ := attributevalue.MarshalMap(i1)
		m2, _ := attributevalue.MarshalMap(i2)
		wrapper.MockQuery(dynamoMock, []map[string]types.AttributeValue{m1, m2}, nil)

		res, err := repo.FindAll(BuildContext(), repo.(*dynamoRepository).wrapper.EmptySpecification())
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("FindAll Error", func(t *testing.T) {
		wrapper.MockQuery(dynamoMock, nil, genericError)
		res, err := repo.FindAll(BuildContext(), repo.(*dynamoRepository).wrapper.EmptySpecification())
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}
