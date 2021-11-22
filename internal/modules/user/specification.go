package user

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/spf13/viper"
)

func WithSMSNotifications() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) (string, expression.Builder) {
		return userTable(), builderExpression.WithFilter(expression.Name("send_daily_sms").Equal(expression.Value(true)).And(withChurchId(ctx)))
	}
}

func WithEmailNotifications() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) (string, expression.Builder) {
		return userTable(), builderExpression.WithFilter(expression.Name("send_weekly_email").Equal(expression.Value(true)).And(withChurchId(ctx)))
	}
}

func WithUserName(username string) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) (string, expression.Builder) {
		userExpr := expression.Name("username").Equal(expression.Value(username))
		return userTable(), builderExpression.WithFilter(userExpr)
	}
}

func withChurchId(ctx context.Context) expression.ConditionBuilder {
	return expression.Name("church_id").Equal(expression.Value(domain.GetChurchID(ctx)))
}

func userTable() string {
	return viper.GetString("tables.user")
}
