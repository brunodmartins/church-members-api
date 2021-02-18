package reports

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"sort"
	"strings"

	"github.com/BrunoDM2943/church-members-api/entity"
	tr "github.com/BrunoDM2943/church-members-api/infra/i18n"
	"github.com/BrunoDM2943/church-members-api/reports/pdf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//go:generate mockgen -source=./reports.go -destination=./mock/reports_mock.go
type ReportsGenerator interface {
	LegalReport() ([]byte, error)
	MemberReport() ([]byte, error)
	BirthdayReport() ([]byte, error)
	MariageReport() ([]byte, error)
	ClassificationReport(classification string) ([]byte, error)
}

type reportService struct {
	repo ReportRepository
}

func NewReportsGenerator(repo ReportRepository) ReportsGenerator {
	return reportService{
		repo,
	}
}

func (report reportService) BirthdayReport() ([]byte, error) {
	members, err := report.repo.FindMembersActive()
	if err != nil {
		return nil, err
	}

	sort.Sort(entity.SortByBirthDay(members))
	data := transformToCSVData(members, func(m *entity.Member) []string {
		return []string{
			m.Person.GetFullName(),
			m.Person.BirthDate.Format("02/01"),
		}
	})
	return writeData(data), nil
}

func writeData(data [][]string) []byte {
	byteArr := &bytes.Buffer{}
	buffer := bufio.NewWriter(byteArr)
	writter := csv.NewWriter(buffer)
	writter.WriteAll(data)
	return byteArr.Bytes()
}

func (report reportService) MariageReport() ([]byte, error) {

	members, err := report.repo.FindMembersActiveAndMarried()

	if err != nil {
		return nil, err
	}

	sort.Sort(entity.SortByMarriageDay(members))
	data := transformToCSVData(members, func(m *entity.Member) []string {
		return []string{
			m.Person.GetFullName() + "&" + m.Person.SpousesName,
			m.Person.MarriageDate.Format("02/01"),
		}
	})

	return writeData(data), nil
}

func (report reportService) MemberReport() ([]byte, error) {
	members, err := report.repo.FindMembersActive()
	if err != nil {
		return nil, err
	}
	sort.Sort(entity.SortByName(members))
	return pdf.BuildPdf(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Reports.Title.Default",
			Other: "Member's report",
		},
	}), members)
}

func (report reportService) ClassificationReport(classification string) ([]byte, error) {
	members, err := report.repo.FindMembersActive()
	if err != nil {
		return nil, err
	}
	members = filterClassification(classification, members)
	sort.Sort(entity.SortByName(members))
	return pdf.BuildPdf(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Reports.Title.Default",
			Other: "Member's report",
		},
	}), members)
}

func filterClassification(classification string, members []*entity.Member) []*entity.Member {
	filtered := []*entity.Member{}
	for _, v := range members {
		if strings.ToLower(v.Classification()) == classification {
			filtered = append(filtered, v)
		}
	}
	return filtered
}


func (report reportService) LegalReport() ([]byte, error) {
	members, err := report.repo.FindMembersActive()
	if err != nil {
		return nil, err
	}
	members = filterChildren(members)
	sort.Sort(entity.SortByName(members))
	return pdf.BuildPdf(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Reports.Title.Legal",
			Other: "Member's report - Legal",
		},
	}), members)
}

func filterChildren(members []*entity.Member) []*entity.Member {
	filtered := []*entity.Member{}
	for _, v := range members {
		if v.Classification() != "Children" {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func transformToCSVData(members []*entity.Member, clojure func(*entity.Member) []string) [][]string {
	data := [][]string{}
	data = append(data, []string{tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Name",
			Other: "Name",
		},
	}), tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Date",
			Other: "Date",
		},
	})})

	for _, member := range members {
		data = append(data, clojure(member))
	}

	return data
}
