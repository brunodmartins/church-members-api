package member

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"time"
)

//QuerySpecificationBuilder allows a client to add dynamic filters to a Query
type QuerySpecificationBuilder struct {
	values map[string]interface{}
}

type QuerySpecification func(builderExpression expression.Builder) expression.Builder

type Specification func(member *domain.Member) bool

func (spec *QuerySpecificationBuilder) AddFilter(key string, value interface{}) {
	if spec.values == nil {
		spec.values = make(map[string]interface{})
	}
	spec.values[key] = value
}

//ToSpecification apply filters to a search on the repo
func (spec *QuerySpecificationBuilder) ToSpecification() QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		var conditions []expression.ConditionBuilder
		if spec.values["gender"] != nil {
			conditions = append(conditions, expression.Name("gender").Equal(expression.Value(spec.values["gender"].(string))))
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

func OnlyActive() QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(activeCondition(true))
	}
}

func OnlyMarriage() QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("marriageDate").AttributeExists().And(activeCondition(true)))
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

func LastMarriages(startDate, endDate time.Time) QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("marriageDateShort").Between(expression.Value(convertDate(startDate)), expression.Value(convertDate(endDate))))
	}
}

func LastBirths(startDate, endDate time.Time) QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("birthDateShort").Between(expression.Value(convertDate(startDate)), expression.Value(convertDate(endDate))))
	}
}

func WithBirthday(date time.Time) QuerySpecification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("birthDateShort").Equal(expression.Value(convertDate(date))))
	}
}

func convertDate(date time.Time) string {
	month := date.Month()
	day := date.Day()
	fmtMonth := fmt.Sprintf("%d", month)
	fmtDay := fmt.Sprintf("%d", day)
	if month < 10 {
		fmtMonth = fmt.Sprintf("0%d", month)
	}
	if day < 10 {
		fmtDay = fmt.Sprintf("0%d", day)
	}
	return fmt.Sprintf("%s-%s", fmtMonth, fmtDay)
}
