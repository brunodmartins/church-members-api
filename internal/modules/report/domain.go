package report

import "github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"

type Report struct {
	Name         string
	Type         reportType.Type
	URL          string
	CreationDate string
}
