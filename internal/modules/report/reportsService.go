package report

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/BrunoDM2943/church-members-api/internal/modules/report/file"
	"github.com/BrunoDM2943/church-members-api/platform/i18n"
	"sort"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
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
	memberService  member.Service
	fileBuilder    file.Builder
	messageService *i18n.MessageService
}

func NewReportService(memberService member.Service, fileBuilder file.Builder) Service {
	return &reportService{
		memberService,
		fileBuilder,
		i18n.GetMessageService(),
	}
}

func (report reportService) BirthdayReport() ([]byte, error) {
	members, err := report.memberService.SearchMembers(member.CreateActiveFilter())
	if err != nil {
		return nil, err
	}

	sort.Sort(domain.SortByBirthDay(members))
	csvOut := file.TransformToCSVData(buildCSVData(members), report.getCSVColumns(), func(row file.Data) []string {
		member := row.Value.(*domain.Member)
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

	members, err := report.memberService.SearchMembers(createMarriageFilter())

	if err != nil {
		return nil, err
	}

	sort.Sort(domain.SortByMarriageDay(members))

	csvOut := file.TransformToCSVData(buildCSVData(members), report.getCSVColumns(), func(row file.Data) []string {
		member := row.Value.(*domain.Member)
		return []string{
			member.Person.GetFullName() + "&" + member.Person.SpousesName,
			member.Person.MarriageDate.Format("02/01"),
		}
	})

	return writeData(csvOut), nil
}

func (report reportService) MemberReport() ([]byte, error) {
	members, err := report.memberService.SearchMembers(member.CreateActiveFilter())
	if err != nil {
		return nil, err
	}
	sort.Sort(domain.SortByName(members))
	return report.fileBuilder.BuildFile(report.messageService.GetMessage("Reports.Title.Default", "Member's report"), members)
}

func (report reportService) ClassificationReport(classification enum.Classification) ([]byte, error) {
	members, err := report.memberService.SearchMembers(member.CreateActiveFilter())
	if err != nil {
		return nil, err
	}
	members = selectByClassification(classification, members)
	sort.Sort(domain.SortByName(members))
	return report.fileBuilder.BuildFile(report.messageService.GetMessage("Reports.Title.Default", "Member's report"), members)
}


func (report reportService) LegalReport() ([]byte, error) {
	members, err := report.memberService.SearchMembers(member.CreateActiveFilter())
	if err != nil {
		return nil, err
	}
	members = filterChildren(members)
	sort.Sort(domain.SortByName(members))
	return report.fileBuilder.BuildFile(report.messageService.GetMessage("Reports.Title.Legal", "Member's report - Legal"), members)
}

func filterChildren(members []*domain.Member) []*domain.Member {
	var filtered []*domain.Member
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

func buildCSVData(members []*domain.Member) []file.Data {
	var data []file.Data
	for _, member := range members {
		data = append(data, file.Data{Value: member})
	}
	return data
}
