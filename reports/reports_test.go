package reports

import (
	"errors"
	"testing"
	"time"

	"github.com/BrunoDM2943/church-members-api/entity"
	mock_repo "github.com/BrunoDM2943/church-members-api/reports/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransformCSVToData(t *testing.T) {
	t1, _ := time.Parse("02/01/2006", "07/06/2020")
	t2, _ := time.Parse("02/01/2006", "22/03/2020")
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName: "Teste",
				LastName:  "Teste",
				BirthDate: t1,
			},
		},
		{
			Person: entity.Person{
				FirstName: "Teste 2",
				LastName:  "Teste 2",
				BirthDate: t2,
			},
		},
	}
	data := transformToCSVData(members, func(m *entity.Member) []string {
		return []string{
			m.Person.GetFullName(),
			m.Person.BirthDate.Format("02/01"),
		}
	})
	assert.Equal(t, 3, len(data))
	assert.Equal(t, []string{"Name", "Date"}, data[0])
	assert.Equal(t, []string{"Teste 2 Teste 2", "22/03"}, data[2])
	assert.Equal(t, []string{"Teste Teste", "07/06"}, data[1])
}

func TestBirthdayReportSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)

	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName: "Teste",
				LastName:  "Teste",
				BirthDate: time.Now(),
			},
		},
		{
			Person: entity.Person{
				FirstName: "Teste 2",
				LastName:  "Teste 2",
				BirthDate: time.Now(),
			},
		},
	}
	repo.EXPECT().FindMembersActive().Return(members, nil)
	out, err := service.BirthdayReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestBirthdayReportSuccessErrorDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)

	repo.EXPECT().FindMembersActive().Return(nil, errors.New("Error"))
	out, err := service.BirthdayReport()
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestMarriageReportSuccessErrorDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)

	repo.EXPECT().FindMembersActiveAndMarried().Return(nil, errors.New("Error"))
	out, err := service.MariageReport()
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestMarriageReportSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)

	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName:    "Esposa",
				LastName:     "Teste",
				MarriageDate: time.Now(),
				SpousesName:  "Marido Teste",
			},
		},
	}
	repo.EXPECT().FindMembersActiveAndMarried().Return(members, nil)
	out, err := service.MariageReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateMemberReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)
	dtNascimento, _ := time.Parse("2006/01/02", "2020/07/06")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    dtNascimento,
				MarriageDate: dtCasamento,
				SpousesName:  "Test spuse",
				Contact: entity.Contact{
					CellPhoneArea: 99,
					CellPhone:     1234567890,
					PhoneArea:     99,
					Phone:         12345678,
					Email:         "teste@test.com",
				},
				Address: entity.Address{
					District: "9",
					City:     "Does not sleep",
					State:    "My-State",
					Address:  "XXXXX",
					Number:   9,
				},
			},
		},
	}
	repo.EXPECT().FindMembersActive().Return(members, nil)
	out, err := service.MemberReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateClassificationReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)
	dtNascimento, _ := time.Parse("2006/01/02", "1990/07/06")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    dtNascimento,
				MarriageDate: dtCasamento,
				SpousesName:  "Test spuse",
				Contact: entity.Contact{
					CellPhoneArea: 99,
					CellPhone:     1234567890,
					PhoneArea:     99,
					Phone:         12345678,
					Email:         "teste@test.com",
				},
				Address: entity.Address{
					District: "9",
					City:     "Does not sleep",
					State:    "My-State",
					Address:  "XXXXX",
					Number:   9,
				},
			},
		},
	}
	repo.EXPECT().FindMembersActive().Return(members, nil)
	out, err := service.ClassificationReport("adult")
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateMemberReportFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)

	repo.EXPECT().FindMembersActive().Return(nil, errors.New("Error"))
	_, err := service.MemberReport()
	assert.NotNil(t, err)
}

func TestGenerateClassificationReportFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)

	repo.EXPECT().FindMembersActive().Return(nil, errors.New("Error"))
	_, err := service.ClassificationReport("adult")
	assert.NotNil(t, err)
}

func TestGenerateLegalReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)
	dtNascimento, _ := time.Parse("2006/01/02", "2020/06/07")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    dtNascimento,
				MarriageDate: dtCasamento,
				SpousesName:  "Test spuse",
				Contact: entity.Contact{
					CellPhoneArea: 99,
					CellPhone:     1234567890,
					PhoneArea:     99,
					Phone:         12345678,
					Email:         "teste@test.com",
				},
				Address: entity.Address{
					District: "9",
					City:     "Does not sleep",
					State:    "My-State",
					Address:  "XXXXX",
					Number:   9,
				},
			},
		},
	}
	repo.EXPECT().FindMembersActive().Return(members, nil)
	out, err := service.LegalReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateLegalReportFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repo.NewMockReportRepository(ctrl)
	service := NewReportsGenerator(repo)

	repo.EXPECT().FindMembersActive().Return(nil, errors.New("Error"))
	_, err := service.LegalReport()
	assert.NotNil(t, err)
}

func TestFilterByClassification(t *testing.T) {
	adult, _ := time.Parse("2006/01/02", "1990/07/06")
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName:    "Adult",
				BirthDate:    adult,
			},
		},
		{
			Person: entity.Person{
				FirstName:    "Children",
				BirthDate:    time.Now(),
			},
		},
	}
	assert.Equal(t, 1, len(filterClassification("adult", members)))
}