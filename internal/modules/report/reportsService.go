package report

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"sort"
	"strings"

	"github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"
	"github.com/brunodmartins/church-members-api/internal/services/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
	GetReport(ctx context.Context, name reportType.Type) (string, error)
	ListReports(ctx context.Context) []Report
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
	fileName, err := reportType.GetFileName(reportType.Type(strings.ToLower(classification.String())))
	if err != nil {
		return err
	}
	return report.storageService.SaveFile(ctx, fileName, result)
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
		i18n.GetMessage(ctx, "Domain.Date"),
	}
}

func (srv *reportService) GetReport(ctx context.Context, report reportType.Type) (string, error) {
	logrus.WithField("church_id", domain.GetChurchID(ctx)).Infof("Getting report %s", report)
	fileName, err := reportType.GetFileName(report)
	if err != nil {
		return "", err
	}
	return srv.storageService.GetFileURL(ctx, fileName)
}

func (srv *reportService) ListReports(ctx context.Context) []Report {
	reports := make([]Report, 0, len(reportType.ReportsTypes))
	creationDates := srv.getReportsCreationDate(ctx)

	for _, report := range reportType.ReportsTypes {
		reports = append(reports, Report{
			Name:         i18n.GetMessage(ctx, "Reports.Title."+cases.Title(language.English).String(string(report))),
			Type:         report,
			URL:          fmt.Sprintf("/reports/%s", report),
			CreationDate: creationDates[report],
		})
	}

	return reports
}

func (srv *reportService) getReportsCreationDate(ctx context.Context) map[reportType.Type]string {
	var result = make(map[reportType.Type]string)
	for _, report := range reportType.ReportsTypes {
		fileName, err := reportType.GetFileName(report)
		if err != nil {
			logrus.WithField("church_id", domain.GetChurchID(ctx)).Errorf("Error parsing last modified date for report %s: %v", report, err)
			result[report] = "Error getting last modified date"
			continue
		}
		metadata, err := srv.storageService.GetFileMetadata(ctx, fileName)
		if err != nil {
			logrus.WithField("church_id", domain.GetChurchID(ctx)).Errorf("Error parsing last modified date for report %s: %v", report, err)
			result[report] = "Error getting last modified date"
			continue
		}
		if metadata.LastModified == nil {
			logrus.WithField("church_id", domain.GetChurchID(ctx)).Errorf("Error parsing last modified date for report %s: missing last modified date", report)
			result[report] = "Error getting last modified date"
			continue
		}
		result[report] = metadata.LastModified.Format("02/01/2006 15:04:05")
	}

	return result
}

func buildCSVData(members []*domain.Member) []file.Data {
	var data []file.Data
	for _, member := range members {
		data = append(data, file.Data{Value: member})
	}
	return data
}
