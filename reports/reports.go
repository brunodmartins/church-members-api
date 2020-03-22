package reports

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
	member "github.com/BrunoDM2943/church-members-api/member/service"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./reports.go -destination=./mock/reports_mock.go
type ReportsGenerator interface {
	//JudicialReport() ([]byte, error)
	//MemberReport() ([]byte, error)
	BirthdayReport() ([]byte, error)
	//MariageReport() ([]byte, error)
}

type ReportService struct {
	memberService member.IMemberService
	repository    ReportRepository
}

func NewReportsGenerator(memberService member.IMemberService, repository ReportRepository) ReportsGenerator {
	return ReportService{
		memberService,
		repository,
	}
}

func (report ReportService) BirthdayReport() ([]byte, error) {
	members, err := report.repository.FindMonthBirthday(time.Now())

	if err != nil {
		return nil, err
	}

	data := transformToCSVData(members)
	byteArr := &bytes.Buffer{}
	buffer := bufio.NewWriter(byteArr)
	writter := csv.NewWriter(buffer)
	writter.WriteAll(data)

	if writter.Error() != nil {
		logrus.Error("Error generating CSV", writter.Error())
		return nil, writter.Error()
	}

	return byteArr.Bytes(), nil
}

func transformToCSVData(members []*entity.Pessoa) [][]string {
	data := [][]string{}
	data = append(data, []string{"Nome", "Data"})

	for _, member := range members {
		data = append(data, []string{
			member.GetFullName(),
			member.DtNascimento.Format("02/01/2006"),
		})
	}

	return data
}
