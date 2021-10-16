package user

import (
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func WithSMSNotifications() wrapper.QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("send_daily_sms").Equal(expression.Value(true)))
	}
}

func WithEmailNotifications() wrapper.QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("send_weekly_email").Equal(expression.Value(true)))
	}
}

func WithUserName(username string) wrapper.QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		userExpr := expression.Name("username").Equal(expression.Value(username))
		return builderExpression.WithFilter(userExpr)
	}
}
