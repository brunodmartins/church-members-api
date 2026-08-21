package reportType

import (
	"fmt"
	"slices"
)

type Type string

const (
	MEMBER    = "members"
	LEGAL     = "legal"
	BIRTHDATE = "birthdate"
	MARRIAGE  = "marriage"
	CHILDREN  = "children"
	TEEN      = "teen"
	YOUNG     = "young"
	ADULT     = "adult"
)

const (
	birthDayReportName = "birthday_report.csv"
	marriageReportName = "marriage_report.csv"
	memberReportName   = "members_report.pdf"
	legalReportName    = "legal_report.pdf"
	childrenReportName = "children_report.pdf"
	teenReportName     = "teen_report.pdf"
	youngReportName    = "young_report.pdf"
	adultReportName    = "adult_report.pdf"
)

var ReportsTypes = []Type{MEMBER, LEGAL, BIRTHDATE, MARRIAGE, CHILDREN, TEEN, YOUNG, ADULT}

func IsValidReport(name Type) bool {
	return slices.Contains(ReportsTypes, name)
}

func GetFileName(reportTypeName Type) (string, error) {
	result := ""
	switch reportTypeName {
	case LEGAL:
		result = legalReportName
	case MEMBER:
		result = memberReportName
	case CHILDREN:
		result = childrenReportName
	case TEEN:
		result = teenReportName
	case YOUNG:
		result = youngReportName
	case ADULT:
		result = adultReportName
	case BIRTHDATE:
		result = birthDayReportName
	case MARRIAGE:
		result = marriageReportName
	}
	if result == "" {
		return "", fmt.Errorf("invalid report type: %s", reportTypeName)
	}
	return result, nil
}
