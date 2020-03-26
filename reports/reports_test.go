package reports

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
	mock_members "github.com/BrunoDM2943/church-members-api/member/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransformCSVToData(t *testing.T) {
	t1, _ := time.Parse("02/01/2006", "10/05/1995")
	t2, _ := time.Parse("02/01/2006", "22/03/2020")
	members := []*entity.Membro{
		&entity.Membro{
			Pessoa: entity.Pessoa{
				Nome:         "Teste",
				Sobrenome:    "Teste",
				DtNascimento: t1,
			},
		},
		&entity.Membro{
			Pessoa: entity.Pessoa{
				Nome:         "Teste 2",
				Sobrenome:    "Teste 2",
				DtNascimento: t2,
			},
		},
	}
	data := transformToCSVData(members, func(m *entity.Membro) []string {
		return []string{
			m.Pessoa.GetFullName(),
			m.Pessoa.DtNascimento.Format("02/01"),
		}
	})
	assert.Equal(t, 3, len(data))
	assert.Equal(t, []string{"Nome", "Data"}, data[0])
	assert.Equal(t, []string{"Teste 2 Teste 2", "22/03"}, data[2])
	assert.Equal(t, []string{"Teste Teste", "10/05"}, data[1])
}

func TestBirthdayReportSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_members.NewMockIMemberService(ctrl)
	service := NewReportsGenerator(memberService)

	members := []*entity.Membro{
		&entity.Membro{
			Pessoa: entity.Pessoa{
				Nome:         "Teste",
				Sobrenome:    "Teste",
				DtNascimento: time.Now(),
			},
		},
		&entity.Membro{
			Pessoa: entity.Pessoa{
				Nome:         "Teste 2",
				Sobrenome:    "Teste 2",
				DtNascimento: time.Now(),
			},
		},
	}
	memberService.EXPECT().FindMembers(gomock.Any()).Return(members, nil)
	out, err := service.BirthdayReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestBirthdayReportSuccessErrorDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_members.NewMockIMemberService(ctrl)
	service := NewReportsGenerator(memberService)

	memberService.EXPECT().FindMembers(gomock.Any()).Return(nil, errors.New("Error"))
	out, err := service.BirthdayReport()
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestMarriageReportSuccessErrorDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_members.NewMockIMemberService(ctrl)
	service := NewReportsGenerator(memberService)

	memberService.EXPECT().FindMembers(gomock.Any()).Return(nil, errors.New("Error"))
	out, err := service.MariageReport()
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestMarriageReportSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_members.NewMockIMemberService(ctrl)
	service := NewReportsGenerator(memberService)

	members := []*entity.Membro{
		&entity.Membro{
			Pessoa: entity.Pessoa{
				Nome:        "Esposa",
				Sobrenome:   "Teste",
				DtCasamento: time.Now(),
				NomeConjuge: "Marido Teste",
			},
		},
	}
	memberService.EXPECT().FindMembers(gomock.Any()).Return(members, nil)
	out, err := service.MariageReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateMemberReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_members.NewMockIMemberService(ctrl)
	service := NewReportsGenerator(memberService)
	dtNascimento, _ := time.Parse("2006/01/02", "1995/05/10")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*entity.Membro{
		&entity.Membro{
			Pessoa: entity.Pessoa{
				Nome:         "Bruno",
				Sobrenome:    "Damasceno Martins",
				DtNascimento: dtNascimento,
				DtCasamento:  dtCasamento,
				NomeConjuge:  "Leticia de Souza Soares da Costa",
				Contato: entity.Contato{
					DDDCelular:  11,
					Celular:     953200587,
					DDDTelefone: 11,
					Telefone:    29435002,
					Email:       "bdm2943@gmail.com",
				},
				Endereco: entity.Endereco{
					Bairro:     "Belem",
					Cidade:     "São Paulo",
					UF:         "SP",
					Logradouro: "Avenida Celso Garcia",
					Numero:     1907,
				},
			},
		},
	}
	memberService.EXPECT().FindMembers(gomock.Any()).Return(members, nil)
	out, err := service.MemberReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
	ioutil.WriteFile("./report.pdf", out, 0644)
}
