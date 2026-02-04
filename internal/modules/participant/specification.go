package participant

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
)

// QueryBuilder allows a client to add dynamic filters to a Query
type QueryBuilder struct {
	values map[string]interface{}
}

type Specification func(participant *domain.Participant) bool

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
		if spec.values["name"] != nil {
			keyCondition = keyCondition.And(expression.Key("name").BeginsWith(spec.values["name"].(string)))
			index = nameIndex()
		}
		if spec.values["active"] != nil {
			filters = append(filters, activeCondition(spec.values["active"].(bool)))
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

func applySpecifications(participants []*domain.Participant, specification []Specification) []*domain.Participant {
	var filtered []*domain.Participant
	for _, participant := range participants {
		allSpecTrue := true
		for _, spec := range specification {
			allSpecTrue = allSpecTrue && spec(participant)
		}
		if allSpecTrue {
			filtered = append(filtered, participant)
		}
	}
	return filtered
}

func withChurchKey(ctx context.Context) expression.KeyConditionBuilder {
	return expression.Key("church_id").Equal(expression.Value(domain.GetChurchID(ctx)))
}

func nameIndex() string {
	return "nameIndex"
}

func activeCondition(value bool) expression.ConditionBuilder {
	return expression.Name("active").Equal(expression.Value(value))
}
