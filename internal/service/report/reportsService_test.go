package report

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/spf13/viper"
	"testing"
	"time"

	mock_repository "github.com/BrunoDM2943/church-members-api/internal/repository/mock"
	mock_file "github.com/BrunoDM2943/church-members-api/internal/storage/file/mock"

	"github.com/BrunoDM2943/church-members-api/internal/constants/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init(){
	viper.Set("bundles.location", "../../../bundles")
}

func TestBirthdayReportSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)
	now := time.Now()
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName: "Teste",
				LastName:  "Teste",
				BirthDate: &now,
			},
		},
		{
			Person: entity.Person{
				FirstName: "Teste 2",
				LastName:  "Teste 2",
				BirthDate: &now,
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
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)

	repo.EXPECT().FindMembersActive().Return(nil, errors.New("Error"))
	out, err := service.BirthdayReport()
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestMarriageReportSuccessErrorDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)

	repo.EXPECT().FindMembersActiveAndMarried().Return(nil, errors.New("Error"))
	out, err := service.MarriageReport()
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestMarriageReportSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)
	now := time.Now()
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName:    "Esposa",
				LastName:     "Teste",
				MarriageDate: &now,
				SpousesName:  "Marido Teste",
			},
		},
	}
	repo.EXPECT().FindMembersActiveAndMarried().Return(members, nil)
	out, err := service.MarriageReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateMemberReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)
	dtNascimento, _ := time.Parse("2006/01/02", "2020/07/06")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    &dtNascimento,
				MarriageDate: &dtCasamento,
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
	fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Eq(members)).Return([]byte{}, nil)
	out, err := service.MemberReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateClassificationReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)
	dtNascimento, _ := time.Parse("2006/01/02", "1990/07/06")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    &dtNascimento,
				MarriageDate: &dtCasamento,
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
	fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Eq(members)).Return([]byte{}, nil)
	out, err := service.ClassificationReport(enum.ADULT)
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateMemberReportFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)

	repo.EXPECT().FindMembersActive().Return(nil, errors.New("Error"))
	_, err := service.MemberReport()
	assert.NotNil(t, err)
}

func TestGenerateClassificationReportFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)

	repo.EXPECT().FindMembersActive().Return(nil, errors.New("Error"))
	_, err := service.ClassificationReport(enum.ADULT)
	assert.NotNil(t, err)
}

func TestGenerateLegalReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)

	dtNascimento, _ := time.Parse("2006/01/02", "2020/06/07")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    &dtNascimento,
				MarriageDate: &dtCasamento,
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
	fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any()).Return([]byte{}, nil)
	out, err := service.LegalReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateLegalReportFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockMemberRepository(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := NewReportService(repo, fileBuilder)

	repo.EXPECT().FindMembersActive().Return(nil, errors.New("Error"))
	_, err := service.LegalReport()
	assert.NotNil(t, err)
}

func TestFilterByClassification(t *testing.T) {
	adult, _ := time.Parse("2006/01/02", "1990/07/06")
	now := time.Now()
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName: "Adult",
				BirthDate: &adult,
			},
		},
		{
			Person: entity.Person{
				FirstName: "Children",
				BirthDate: &now,
			},
		},
	}
	assert.Equal(t, 1, len(filterClassification(enum.ADULT, members)))
}

func TestFilterByClassificationRemovingChildren(t *testing.T) {
	adult, _ := time.Parse("2006/01/02", "1990/07/06")
	now := time.Now()
	members := []*entity.Member{
		{
			Person: entity.Person{
				FirstName: "Adult",
				BirthDate: &adult,
			},
		},
		{
			Person: entity.Person{
				FirstName: "Children",
				BirthDate: &now,
			},
		},
	}
	assert.Equal(t, 1, len(filterChildren(members)))
}

