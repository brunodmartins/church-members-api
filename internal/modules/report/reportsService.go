package report

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum/reportType"
	"github.com/BrunoDM2943/church-members-api/internal/services/storage"
	"github.com/sirupsen/logrus"
	"sort"

	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/BrunoDM2943/church-members-api/internal/modules/report/file"
	"github.com/BrunoDM2943/church-members-api/platform/i18n"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
)

//go:generate mockgen -source=./reportsService.go -destination=./mock/reports_mock.go
type Service interface {
	LegalReport(ctx context.Context) ([]byte, error)
	MemberReport(ctx context.Context) ([]byte, error)
	BirthdayReport(ctx context.Context) ([]byte, error)
	MarriageReport(ctx context.Context) ([]byte, error)
	ClassificationReport(ctx context.Context, classification enum.Classification) ([]byte, error)

	GetReport(ctx context.Context, name string) (string, error)
}

type reportService struct {
	memberService  member.Service
	fileBuilder    file.Builder
	messageService *i18n.MessageService
	storageService storage.Service
}

const (
	birthDayReportName       = "birthday_report.csv"
	marriageReportName       = "marriage_report.csv"
	memberReportName         = "members_report.pdf"
	classificationReportName = "classification_report.pdf"
	legalReportName          = "legal_report.pdf"
)

func NewReportService(memberService member.Service, fileBuilder file.Builder, storageService storage.Service) Service {
	return &reportService{
		memberService,
		fileBuilder,
		i18n.GetMessageService(),
		storageService,
	}
}

func (report reportService) BirthdayReport(ctx context.Context) ([]byte, error) {
	members, err := report.memberService.SearchMembers(ctx, member.OnlyActive())
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
	result := writeData(csvOut)
	return result, report.storageService.SaveFile(ctx, birthDayReportName, result)
}

func writeData(data [][]string) []byte {
	byteArr := &bytes.Buffer{}
	buffer := bufio.NewWriter(byteArr)
	writter := csv.NewWriter(buffer)
	writter.WriteAll(data)
	return byteArr.Bytes()
}

func (report reportService) MarriageReport(ctx context.Context) ([]byte, error) {

	members, err := report.memberService.SearchMembers(ctx, member.OnlyMarriage())

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

	result := writeData(csvOut)
	return result, report.storageService.SaveFile(ctx, marriageReportName, result)
}

func (report reportService) MemberReport(ctx context.Context) ([]byte, error) {
	members, err := report.memberService.SearchMembers(ctx, member.OnlyActive())
	if err != nil {
		return nil, err
	}
	sort.Sort(domain.SortByName(members))
	result, err := report.fileBuilder.BuildFile(report.messageService.GetMessage("Reports.Title.Default", "Member's report"), domain.GetChurch(ctx), members)
	if err != nil {
		return nil, err
	}
	return result, report.storageService.SaveFile(ctx, memberReportName, result)
}

func (report reportService) ClassificationReport(ctx context.Context, classification enum.Classification) ([]byte, error) {
	members, err := report.memberService.SearchMembers(ctx, member.OnlyActive(), member.OnlyByClassification(classification))
	if err != nil {
		return nil, err
	}
	sort.Sort(domain.SortByName(members))
	result, err := report.fileBuilder.BuildFile(report.messageService.GetMessage("Reports.Title.Default", "Member's report"), domain.GetChurch(ctx), members)
	if err != nil {
		return nil, err
	}
	return result, report.storageService.SaveFile(ctx, classificationReportName, result)
}

func (report reportService) LegalReport(ctx context.Context) ([]byte, error) {
	members, err := report.memberService.SearchMembers(ctx, member.OnlyActive(), member.OnlyLegalMembers())
	if err != nil {
		return nil, err
	}
	sort.Sort(domain.SortByName(members))
	result, err := report.fileBuilder.BuildFile(report.messageService.GetMessage("Reports.Title.Legal", "Member's report - Legal"), domain.GetChurch(ctx), members)
	if err != nil {
		return nil, err
	}
	return result, report.storageService.SaveFile(ctx, legalReportName, result)
}

func (report *reportService) getCSVColumns() []string {
	return []string{
		report.messageService.GetMessage("Domain.Name", "Name"),
		report.messageService.GetMessage("Domain.Date", "Date"),
	}
}

func (report reportService) GetReport(ctx context.Context, reportType string) (string, error) {
	logrus.WithField("church_id", domain.GetChurchID(ctx)).Infof("Getting report %s", reportType)
	fileName, err := getFileName(reportType)
	if err != nil {
		return "", err
	}
	return report.storageService.GetFileURL(ctx, fileName)
}

func getFileName(reportTypeName string) (string, error) {
	result := ""
	switch reportTypeName {
	case reportType.LEGAL:
		result = legalReportName
	case reportType.MEMBER:
		result = memberReportName
	case reportType.CLASSIFICATION:
		result = classificationReportName
	case reportType.BIRTHDATE:
		result = birthDayReportName
	case reportType.MARRIAGE:
		result = marriageReportName
	}
	if result == "" {
		return "", errors.New("invalid report type: " + reportTypeName)
	}
	return result, nil
}

func buildCSVData(members []*domain.Member) []file.Data {
	var data []file.Data
	for _, member := range members {
		data = append(data, file.Data{Value: member})
	}
	return data
}
