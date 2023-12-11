package member

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/brunodmartins/church-members-api/platform/utils"
	"time"
)

// QueryBuilder allows a client to add dynamic filters to a Query
type QueryBuilder struct {
	values map[string]interface{}
}

type Specification func(member *domain.Member) bool

func (spec *QueryBuilder) AddFilter(key string, value interface{}) {
	if spec.values == nil {
		spec.values = make(map[string]interface{})
	}
	spec.values[key] = value
}

// ToSpecification apply filters to a search on the repo
func (spec *QueryBuilder) ToSpecification() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		var filters []expression.ConditionBuilder
		index := ""
		keyCondition := withChurchKey(ctx)
		if spec.values["gender"] != nil {
			filters = append(filters, expression.Name("gender").Equal(expression.Value(spec.values["gender"].(string))))
		}
		if spec.values["active"] != nil {
			filters = append(filters, activeCondition(spec.values["active"].(bool)))
		}
		if spec.values["name"] != nil {
			keyCondition = keyCondition.And(expression.Key("name").BeginsWith(spec.values["name"].(string)))
			index = nameIndex()
		}
		if len(filters) != 0 {
			builderExpression = builderExpression.WithKeyCondition(keyCondition).WithFilter(spec.mergeFilters(filters))
		} else {
			builderExpression = builderExpression.WithKeyCondition(keyCondition)
		}
		return wrapper.ExpressionBuilder{
			Index:   index,
			Builder: builderExpression,
		}
	}

}

func (spec *QueryBuilder) mergeFilters(filters []expression.ConditionBuilder) expression.ConditionBuilder {
	var finalFilter expression.ConditionBuilder
	for index, filter := range filters {
		if index == 0 {
			finalFilter = filter
		} else {
			finalFilter = finalFilter.And(filter)
		}
	}
	return finalFilter
}

func OnlyActive() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		keyCondition := withChurchKey(ctx)
		return wrapper.ExpressionBuilder{
			Builder: builderExpression.WithKeyCondition(keyCondition).WithFilter(activeCondition(true)),
		}
	}
}

func OnlyInactive() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		keyCondition := withChurchKey(ctx)
		return wrapper.ExpressionBuilder{
			Builder: builderExpression.WithKeyCondition(keyCondition).WithFilter(activeCondition(false)),
		}
	}
}

func OnlyMarriage() wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		keyCondition := withChurchKey(ctx).And(expression.Key("maritalStatus").Equal(expression.Value("MARRIED")))
		return wrapper.ExpressionBuilder{
			Index:   maritalStatusIndex(),
			Builder: builderExpression.WithKeyCondition(keyCondition).WithFilter(activeCondition(true)),
		}
	}
}

func activeCondition(value bool) expression.ConditionBuilder {
	return expression.Name("active").Equal(expression.Value(value))
}

func OnlyLegalMembers() Specification {
	return func(member *domain.Member) bool {
		return member.IsLegal()
	}
}

func OnlyMembershipEndCurrentYear() Specification {
	return func(member *domain.Member) bool {
		return member.MembershipEndCurrentYear()
	}
}

func OnlyByClassification(value enum.Classification) Specification {
	return func(member *domain.Member) bool {
		return member.Classification() == value
	}
}

func applySpecifications(members []*domain.Member, specification []Specification) []*domain.Member {
	var filtered []*domain.Member
	for _, member := range members {
		allSpecTrue := true
		for _, spec := range specification {
			allSpecTrue = allSpecTrue && spec(member)
		}
		if allSpecTrue {
			filtered = append(filtered, member)
		}
	}
	return filtered
}

func LastMarriages(startDate, endDate time.Time) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		maritalStatus := expression.Key("maritalStatus").Equal(expression.Value("MARRIED"))
		builderExpression = builderExpression.WithKeyCondition(maritalStatus.And(withChurchKey(ctx)))
		return wrapper.ExpressionBuilder{
			Index:   maritalStatusIndex(),
			Builder: builderExpression.WithFilter(expression.Name("marriageDateShort").Between(expression.Value(utils.ConvertDate(startDate)), expression.Value(utils.ConvertDate(endDate))).And(activeCondition(true))),
		}
	}
}

func LastBirths(startDate, endDate time.Time) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		dateKey := expression.Key("birthDateShort").Between(expression.Value(utils.ConvertDate(startDate)), expression.Value(utils.ConvertDate(endDate)))
		key := withChurchKey(ctx).And(dateKey)
		return wrapper.ExpressionBuilder{
			Index:   birthDateIndex(),
			Builder: builderExpression.WithKeyCondition(key).WithFilter(activeCondition(true)),
		}
	}
}

func WithBirthday(date time.Time) wrapper.QuerySpecification {
	return func(ctx context.Context, builderExpression expression.Builder) wrapper.ExpressionBuilder {
		dateKey := expression.Key("birthDateShort").Equal(expression.Value(utils.ConvertDate(date)))
		builderExpression = builderExpression.WithKeyCondition(dateKey.And(withChurchKey(ctx)))
		return wrapper.ExpressionBuilder{
			Index:   birthDateIndex(),
			Builder: builderExpression.WithFilter(activeCondition(true)),
		}
	}
}

func withChurchKey(ctx context.Context) expression.KeyConditionBuilder {
	return expression.Key("church_id").Equal(expression.Value(domain.GetChurchID(ctx)))
}

func nameIndex() string {
	return "nameIndex"
}

func birthDateIndex() string {
	return "birthDateIndex"
}

func maritalStatusIndex() string {
	return "maritalStatusIndex"
}
