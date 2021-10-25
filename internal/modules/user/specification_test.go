package user

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithSMSNotifications(t *testing.T) {
	builder := expression.NewBuilder()
	spec := WithSMSNotifications()
	builder = spec(buildContext(), builder)
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 1)
}

func TestWithEmailNotifications(t *testing.T) {
	builder := expression.NewBuilder()
	spec := WithEmailNotifications()
	builder = spec(buildContext(), builder)
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 1)
}

func TestWithUserName(t *testing.T) {
	builder := expression.NewBuilder()
	spec := WithUserName("test")
	builder = spec(buildContext(), builder)
	expression, err := builder.Build()
	assert.Nil(t, err)
	assert.Len(t, expression.Names(), 1)
}
