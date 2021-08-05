package report

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func createMarriageFilter() member.Specification {
	return func(builderExpression expression.Builder) expression.Builder {
		return builderExpression.WithFilter(expression.Name("marriageDate").AttributeExists().And(activeCondition(true)))
	}
}

func activeCondition(value bool) expression.ConditionBuilder {
	return expression.Name("active").Equal(expression.Value(value))
}

func selectByClassification(classification enum.Classification, members []*domain.Member) []*domain.Member {
	var filtered []*domain.Member
	for _, v := range members {
		if v.Classification() == classification {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func selectByNotInClassification(classification enum.Classification, members []*domain.Member) []*domain.Member {
	var filtered []*domain.Member
	for _, v := range members {
		if v.Classification() != classification {
			filtered = append(filtered, v)
		}
	}
	return filtered
}