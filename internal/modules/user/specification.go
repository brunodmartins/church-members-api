package user

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func WithSMSNotifications() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		return wrapper.ExpressionBuilder{
			Builder: builderExpression.WithKeyCondition(withChurchId(ctx)).WithFilter(expression.Name("send_daily_sms").Equal(expression.Value(true))),
		}
	}
}

func WithEmailNotifications() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		return wrapper.ExpressionBuilder{
			Builder: builderExpression.WithKeyCondition(withChurchId(ctx)).WithFilter(expression.Name("send_weekly_email").Equal(expression.Value(true))),
		}
	}
}

func WithUserName(username string) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		userKey := expression.Key("username").Equal(expression.Value(username))
		return wrapper.ExpressionBuilder{
			Builder: builderExpression.WithKeyCondition(withChurchId(ctx).And(userKey)),
		}
	}
}

func withChurchId(ctx context.Context) expression.KeyConditionBuilder {
	return expression.Key("church_id").Equal(expression.Value(domain.GetChurchID(ctx)))
}
