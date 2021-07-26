package report

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/BrunoDM2943/church-members-api/internal/infra/i18n"
	"github.com/BrunoDM2943/church-members-api/internal/repository"
	"github.com/BrunoDM2943/church-members-api/internal/storage/file"
	"sort"

	"github.com/BrunoDM2943/church-members-api/internal/constants/entity"
)

//go:generate mockgen -source=./reportsService.go -destination=./mock/reports_mock.go
type Service interface {
	LegalReport() ([]byte, error)
	MemberReport() ([]byte, error)
	BirthdayReport() ([]byte, error)
	MarriageReport() ([]byte, error)
	ClassificationReport(classification enum.Classification) ([]byte, error)
}

type reportService struct {
	repo        repository.MemberRepository
	fileBuilder file.Builder
	messageService *i18n.MessageService
}

func NewReportService(repo repository.MemberRepository, fileBuilder file.Builder) Service {
	return &reportService{
		repo,
		fileBuilder,
		i18n.GetMessageService(),
	}
}

func (report reportService) BirthdayReport() ([]byte, error) {
	members, err := report.repo.FindMembersActive()
	if err != nil {
		return nil, err
	}

	sort.Sort(entity.SortByBirthDay(members))
	csvOut := file.TransformToCSVData(buildCSVData(members), report.getCSVColumns(), func(row file.Data) []string {
		member := row.Value.(*entity.Member)
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

	sort.Sort(entity.SortByMarriageDay(members))

	csvOut := file.TransformToCSVData(buildCSVData(members), report.getCSVColumns(), func(row file.Data) []string {
		member := row.Value.(*entity.Member)
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
	sort.Sort(entity.SortByName(members))
	return report.fileBuilder.BuildFile(report.messageService.GetMessage("Reports.Title.Default", "Member's report"), members)
}

func (report reportService) ClassificationReport(classification enum.Classification) ([]byte, error) {
	members, err := report.repo.FindMembersActive()
	if err != nil {
		return nil, err
	}
	members = filterClassification(classification, members)
	sort.Sort(entity.SortByName(members))
	return report.fileBuilder.BuildFile(report.messageService.GetMessage("Reports.Title.Default", "Member's report"), members)
}

func filterClassification(classification enum.Classification, members []*entity.Member) []*entity.Member {
	filtered := []*entity.Member{}
	for _, v := range members {
		if v.Classification() == classification {
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
	return report.fileBuilder.BuildFile(report.messageService.GetMessage("Reports.Title.Legal", "Member's report - Legal"), members)
}

func filterChildren(members []*entity.Member) []*entity.Member {
	var filtered []*entity.Member
	for _, v := range members {
		if v.Classification() != enum.CHILDREN {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func (report *reportService) getCSVColumns() []string {
	return []string{
		report.messageService.GetMessage("Domain.Name", "Name"),
		report.messageService.GetMessage("Domain.Date", "Date"),
	}
}

func buildCSVData(members []*entity.Member) []file.Data {
	var data []file.Data
	for _, member := range members {
		data = append(data, file.Data{Value: member})
	}
	return data
}
