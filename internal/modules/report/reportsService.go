package report

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"
	"github.com/brunodmartins/church-members-api/internal/services/storage"
	"github.com/sirupsen/logrus"
	"sort"

	"github.com/brunodmartins/church-members-api/internal/constants/enum"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/modules/report/file"
	"github.com/brunodmartins/church-members-api/platform/i18n"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

//go:generate mockgen -source=./reportsService.go -destination=./mock/reports_mock.go
type Service interface {
	LegalReport(ctx context.Context) error
	MemberReport(ctx context.Context) error
	BirthdayReport(ctx context.Context) error
	MarriageReport(ctx context.Context) error
	ClassificationReport(ctx context.Context, classification enum.Classification) error
	GetReport(ctx context.Context, name string) (string, error)
}

type reportService struct {
	memberService  member.Service
	fileBuilder    file.Builder
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
		storageService,
	}
}

func (report reportService) BirthdayReport(ctx context.Context) error {
	members, err := report.memberService.SearchMembers(ctx, member.OnlyActive())
	if err != nil {
		return err
	}

	sort.Sort(domain.SortByBirthDay(members))
	csvOut := file.TransformToCSVData(buildCSVData(members), report.getCSVColumns(ctx), func(row file.Data) []string {
		member := row.Value.(*domain.Member)
		return []string{
			member.Person.GetFullName(),
			member.Person.BirthDate.Format("02/01"),
		}
	})
	return report.storageService.SaveFile(ctx, birthDayReportName, writeData(csvOut))
}

func writeData(data [][]string) []byte {
	byteArr := &bytes.Buffer{}
	buffer := bufio.NewWriter(byteArr)
	writter := csv.NewWriter(buffer)
	writter.WriteAll(data)
	return byteArr.Bytes()
}

func (report reportService) MarriageReport(ctx context.Context) error {

	members, err := report.memberService.SearchMembers(ctx, member.OnlyMarriage())

	if err != nil {
		return err
	}

	sort.Sort(domain.SortByMarriageDay(members))

	csvOut := file.TransformToCSVData(buildCSVData(members), report.getCSVColumns(ctx), func(row file.Data) []string {
		member := row.Value.(*domain.Member)
		return []string{
			member.Person.GetFullName() + "&" + member.Person.SpousesName,
			member.Person.MarriageDate.Format("02/01"),
		}
	})

	return report.storageService.SaveFile(ctx, marriageReportName, writeData(csvOut))
}

func (report reportService) MemberReport(ctx context.Context) error {
	members, err := report.memberService.SearchMembers(ctx, member.OnlyActive())
	if err != nil {
		return err
	}
	sort.Sort(domain.SortByName(members))
	result, err := report.fileBuilder.BuildFile(ctx, i18n.GetMessage(ctx, "Reports.Title.Default"), domain.GetChurch(ctx), members)
	if err != nil {
		return err
	}
	return report.storageService.SaveFile(ctx, memberReportName, result)
}

func (report reportService) ClassificationReport(ctx context.Context, classification enum.Classification) error {
	members, err := report.memberService.SearchMembers(ctx, member.OnlyActive(), member.OnlyByClassification(classification))
	if err != nil {
		return err
	}
	sort.Sort(domain.SortByName(members))
	result, err := report.fileBuilder.BuildFile(ctx, i18n.GetMessage(ctx, "Reports.Title.Default"), domain.GetChurch(ctx), members)
	if err != nil {
		return err
	}
	return report.storageService.SaveFile(ctx, classificationReportName, result)
}

func (report reportService) LegalReport(ctx context.Context) error {
	members, err := report.memberService.SearchMembers(ctx, member.OnlyActive(), member.OnlyLegalMembers())
	if err != nil {
		return err
	}
	inactiveMembers, err := report.memberService.SearchMembers(ctx, member.OnlyInactive(), member.OnlyLegalMembers(), member.OnlyMembershipEndCurrentYear())
	if err != nil {
		return err
	}
	members = append(members, inactiveMembers...)
	sort.Sort(domain.SortByName(members))
	result, err := report.fileBuilder.BuildFile(ctx, i18n.GetMessage(ctx, "Reports.Title.Legal"), domain.GetChurch(ctx), members)
	if err != nil {
		return err
	}
	return report.storageService.SaveFile(ctx, legalReportName, result)
}

func (report *reportService) getCSVColumns(ctx context.Context) []string {
	return []string{
		i18n.GetMessage(ctx, "Domain.Name"),
		i18n.GetMessage(ctx, "Domain.RetireDate"),
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
