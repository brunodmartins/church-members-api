package user

import (
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func AllowSMSNotifications() wrapper.QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("send_daily_sms").Equal(expression.Value(true)))
	}
}

func AllowEmailNotifications() wrapper.QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("send_weekly_email").Equal(expression.Value(true)))
	}
}
