package user

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWithSMSNotifications(t *testing.T) {
	spec := WithSMSNotifications()
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 2)
}

func TestWithEmailNotifications(t *testing.T) {
	spec := WithEmailNotifications()
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 2)
}

func TestWithUserName(t *testing.T) {
	spec := WithUserName("test")
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 2)
}

func TestWithID(t *testing.T) {
	spec := WithId(uuid.NewString())
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 2)
}

func TestAllUsers(t *testing.T) {
	spec := AllUsers()
	builder := spec(BuildContext(), expression.NewBuilder())
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 1)
	attribute := expression.Values()[":0"]
	assert.Equal(t, "church_id_test", attribute.(*types.AttributeValueMemberS).Value)
}
