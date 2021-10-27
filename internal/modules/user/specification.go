package user

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func WithSMSNotifications() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("send_daily_sms").Equal(expression.Value(true)).And(withChurchId(ctx)))
	}
}

func WithEmailNotifications() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("send_weekly_email").Equal(expression.Value(true)).And(withChurchId(ctx)))
	}
}

func WithUserName(username string) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) expression.Builder {
		userExpr := expression.Name("username").Equal(expression.Value(username))
		return builderExpression.WithFilter(userExpr.And(withChurchId(ctx)))
	}
}

func withChurchId(ctx context.Context) expression.ConditionBuilder {
	return expression.Name("church_id").Equal(expression.Value(domain.GetChurchID(ctx)))
}
