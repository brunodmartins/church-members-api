package reportType

import (
	"errors"
	"slices"
)

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

var validReports = []string{MEMBER, LEGAL, BIRTHDATE, MARRIAGE, CHILDREN, TEEN, YOUNG, ADULT}

func IsValidReport(name string) bool {
	return slices.Contains(validReports, name)
}

func GetFileName(reportTypeName string) (string, error) {
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
		return "", errors.New("invalid report type: " + reportTypeName)
	}
	return result, nil
}
