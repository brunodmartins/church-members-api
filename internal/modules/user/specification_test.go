package user

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
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
