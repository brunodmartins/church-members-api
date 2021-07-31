package member

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

//QuerySpecification allows a client to add dynamic filters to a Query
type QuerySpecification struct {
	values map[string]interface{}
}

type Specification func(builderExpression expression.Builder) expression.Builder

func (spec *QuerySpecification) AddFilter(key string, value interface{}) {
	if spec.values == nil {
		spec.values = make(map[string]interface{})
	}
	spec.values[key] = value
}

//ToSpecification apply filters to a search on the repo
func (spec *QuerySpecification) ToSpecification() Specification {
	return func(builderExpression expression.Builder) expression.Builder {
		var conditions []expression.ConditionBuilder
		if spec.values["person.gender"] != nil {
			conditions = append(conditions, expression.Name("gender").Equal(expression.Value(spec.values["person.gender"].(string))))
		}
		if spec.values["active"] != nil {
			conditions = append(conditions, activeCondition(spec.values["active"].(bool)))
		}
		if spec.values["name"] != nil {
			conditions = append(conditions, expression.Name("name").Contains(spec.values["name"].(string)))
		}
		conditionsSize := len(conditions)
		switch conditionsSize {
		case 0:
			return builderExpression
		case 1:
			return builderExpression.WithFilter(conditions[0])
		case 2:
			return builderExpression.WithFilter(conditions[0].And(conditions[1]))
		default:
			return builderExpression.WithFilter(conditions[0].And(conditions[1], conditions[2:]...))
		}
	}

}

func CreateMarriageFilter() Specification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("marriageDate").AttributeExists().And(activeCondition(true)))
	}
}

func CreateActiveFilter() Specification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(activeCondition(true))
	}
}

func activeCondition(value bool) expression.ConditionBuilder {
	return expression.Name("active").Equal(expression.Value(value))
}