package report

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"sort"
	"strings"

	"github.com/BrunoDM2943/church-members-api/internal/repository"
	"github.com/BrunoDM2943/church-members-api/internal/storage/file"

	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	tr "github.com/BrunoDM2943/church-members-api/internal/infra/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//go:generate mockgen -source=./reportsService.go -destination=./mock/reports_mock.go
type Service interface {
	LegalReport() ([]byte, error)
	MemberReport() ([]byte, error)
	BirthdayReport() ([]byte, error)
	MarriageReport() ([]byte, error)
	ClassificationReport(classification string) ([]byte, error)
}

type reportService struct {
	repo        repository.MemberRepository
	fileBuilder file.Builder
}

func NewReportService(repo repository.MemberRepository, fileBuilder file.Builder) Service {
	return &reportService{
		repo,
		fileBuilder,
	}
}

func (report reportService) BirthdayReport() ([]byte, error) {
	members, err := report.repo.FindMembersActive()
	if err != nil {
		return nil, err
	}

	sort.Sort(model.SortByBirthDay(members))
	csvOut := file.TransformToCSVData(buildCSVData(members), getCSVColumns(), func(row file.Data) []string {
		member := row.Value.(*model.Member)
		return []string{
			member.Person.GetFullName(),
			member.Person.BirthDate.Format("02/01"),
		}
	})
	return writeData(csvOut), nil
}

func writeData(data [][]string) []byte {
	byteArr := &bytes.Buffer{}
	buffer := bufio.NewWriter(byteArr)
	writter := csv.NewWriter(buffer)
	writter.WriteAll(data)
	return byteArr.Bytes()
}

func (report reportService) MarriageReport() ([]byte, error) {

	members, err := report.repo.FindMembersActiveAndMarried()

	if err != nil {
		return nil, err
	}

	sort.Sort(model.SortByMarriageDay(members))

	csvOut := file.TransformToCSVData(buildCSVData(members), getCSVColumns(), func(row file.Data) []string {
		member := row.Value.(*model.Member)
		return []string{
			member.Person.GetFullName() + "&" + member.Person.SpousesName,
			member.Person.MarriageDate.Format("02/01"),
		}
	})

	return writeData(csvOut), nil
}

func (report reportService) MemberReport() ([]byte, error) {
	members, err := report.repo.FindMembersActive()
	if err != nil {
		return nil, err
	}
	sort.Sort(model.SortByName(members))
	return report.fileBuilder.BuildFile(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
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
	sort.Sort(model.SortByName(members))
	return report.fileBuilder.BuildFile(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Reports.Title.Default",
			Other: "Member's report",
		},
	}), members)
}

func filterClassification(classification string, members []*model.Member) []*model.Member {
	filtered := []*model.Member{}
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
	sort.Sort(model.SortByName(members))
	return report.fileBuilder.BuildFile(tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Reports.Title.Legal",
			Other: "Member's report - Legal",
		},
	}), members)
}

func filterChildren(members []*model.Member) []*model.Member {
	filtered := []*model.Member{}
	for _, v := range members {
		if v.Classification() != "Children" {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func getCSVColumns() []string {
	return []string{tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Name",
			Other: "Name",
		},
	}), tr.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Domain.Date",
			Other: "Date",
		},
	})}
}

func buildCSVData(members []*model.Member) []file.Data {
	var data []file.Data
	for _, member := range members {
		data = append(data, file.Data{Value: member})
	}
	return data
}
