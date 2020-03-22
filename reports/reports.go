package reports

import (
	"bufio"
	"bytes"
	"encoding/csv"

	"github.com/BrunoDM2943/church-members-api/entity"
	member "github.com/BrunoDM2943/church-members-api/member/service"
)

//go:generate mockgen -source=./reports.go -destination=./mock/reports_mock.go
type ReportsGenerator interface {
	//JudicialReport() ([]byte, error)
	//MemberReport() ([]byte, error)
	BirthdayReport() ([]byte, error)
	//MariageReport() ([]byte, error)
}

type reportService struct {
	memberService member.IMemberService
}

func NewReportsGenerator(memberService member.IMemberService) ReportsGenerator {
	return reportService{
		memberService,
	}
}

func (report reportService) BirthdayReport() ([]byte, error) {

	members, err := report.memberService.FindMembers(map[string]interface{}{
		"active": true,
	})

	if err != nil {
		return nil, err
	}

	data := transformToCSVData(members)
	byteArr := &bytes.Buffer{}
	buffer := bufio.NewWriter(byteArr)
	writter := csv.NewWriter(buffer)
	writter.WriteAll(data)
	return byteArr.Bytes(), nil
}

func transformToCSVData(members []*entity.Membro) [][]string {
	data := [][]string{}
	data = append(data, []string{"Nome", "Data"})

	for _, member := range members {
		data = append(data, []string{
			member.Pessoa.GetFullName(),
			member.Pessoa.DtNascimento.Format("02/01/2006"),
		})
	}

	return data
}
