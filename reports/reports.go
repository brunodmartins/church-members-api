package reports

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"sort"

	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/reports/pdf"
)

//go:generate mockgen -source=./reports.go -destination=./mock/reports_mock.go
type ReportsGenerator interface {
	//JudicialReport() ([]byte, error)
	MemberReport() ([]byte, error)
	BirthdayReport() ([]byte, error)
	MariageReport() ([]byte, error)
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
	data := transformToCSVData(members, func(m *entity.Membro) []string {
		return []string{
			m.Pessoa.GetFullName(),
			m.Pessoa.DtNascimento.Format("02/01"),
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
	data := transformToCSVData(members, func(m *entity.Membro) []string {
		return []string{
			m.Pessoa.GetFullName() + "&" + m.Pessoa.NomeConjuge,
			m.Pessoa.DtCasamento.Format("02/01"),
		}
	})

	return writeData(data), nil
}

func (report reportService) MemberReport() ([]byte, error) {
	members, _ := report.repo.FindMembersActive()
	sort.Sort(entity.SortByName(members))
	return pdf.BuildPdf("Relatório de Membros", members)

}

func transformToCSVData(members []*entity.Membro, clojure func(*entity.Membro) []string) [][]string {
	data := [][]string{}
	data = append(data, []string{"Nome", "Data"})

	for _, member := range members {
		data = append(data, clojure(member))
	}

	return data
}
