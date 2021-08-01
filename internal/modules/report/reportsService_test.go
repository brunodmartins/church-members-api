package report

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member/mock"
	mock_file2 "github.com/BrunoDM2943/church-members-api/internal/modules/report/file/mock"
	"github.com/spf13/viper"
	"testing"
	"time"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init(){
	viper.Set("bundles.location", "../../../bundles")
}

func TestBirthdayReportSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)
	now := time.Now()
	members := []*domain.Member{
		{
			Person: domain.Person{
				FirstName: "Teste",
				LastName:  "Teste",
				BirthDate: &now,
			},
		},
		{
			Person: domain.Person{
				FirstName: "Teste 2",
				LastName:  "Teste 2",
				BirthDate: &now,
			},
		},
	}
	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(members, nil)
	out, err := service.BirthdayReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestBirthdayReportSuccessErrorDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)

	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(nil, errors.New("error"))
	out, err := service.BirthdayReport()
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestMarriageReportSuccessErrorDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)

	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(nil, errors.New("error"))
	out, err := service.MarriageReport()
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestMarriageReportSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)
	now := time.Now()
	members := []*domain.Member{
		{
			Person: domain.Person{
				FirstName:    "Esposa",
				LastName:     "Teste",
				MarriageDate: &now,
				SpousesName:  "Marido Teste",
			},
		},
	}
	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(members, nil)
	out, err := service.MarriageReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateMemberReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)
	dtNascimento, _ := time.Parse("2006/01/02", "2020/07/06")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*domain.Member{
		{
			Person: domain.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    &dtNascimento,
				MarriageDate: &dtCasamento,
				SpousesName:  "Test spuse",
				Contact: domain.Contact{
					CellPhoneArea: 99,
					CellPhone:     1234567890,
					PhoneArea:     99,
					Phone:         12345678,
					Email:         "teste@test.com",
				},
				Address: domain.Address{
					District: "9",
					City:     "Does not sleep",
					State:    "My-State",
					Address:  "XXXXX",
					Number:   9,
				},
			},
		},
	}
	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(members, nil)
	fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Eq(members)).Return([]byte{}, nil)
	out, err := service.MemberReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateClassificationReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)
	dtNascimento, _ := time.Parse("2006/01/02", "1990/07/06")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*domain.Member{
		{
			Person: domain.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    &dtNascimento,
				MarriageDate: &dtCasamento,
				SpousesName:  "Test spuse",
				Contact: domain.Contact{
					CellPhoneArea: 99,
					CellPhone:     1234567890,
					PhoneArea:     99,
					Phone:         12345678,
					Email:         "teste@test.com",
				},
				Address: domain.Address{
					District: "9",
					City:     "Does not sleep",
					State:    "My-State",
					Address:  "XXXXX",
					Number:   9,
				},
			},
		},
	}
	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(members, nil)
	fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Eq(members)).Return([]byte{}, nil)
	out, err := service.ClassificationReport(enum.ADULT)
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateMemberReportFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)

	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(nil, errors.New("error"))
	_, err := service.MemberReport()
	assert.NotNil(t, err)
}

func TestGenerateClassificationReportFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)

	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(nil, errors.New("error"))
	_, err := service.ClassificationReport(enum.ADULT)
	assert.NotNil(t, err)
}

func TestGenerateLegalReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)

	dtNascimento, _ := time.Parse("2006/01/02", "2020/06/07")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	members := []*domain.Member{
		{
			Person: domain.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    &dtNascimento,
				MarriageDate: &dtCasamento,
				SpousesName:  "Test spuse",
				Contact: domain.Contact{
					CellPhoneArea: 99,
					CellPhone:     1234567890,
					PhoneArea:     99,
					Phone:         12345678,
					Email:         "teste@test.com",
				},
				Address: domain.Address{
					District: "9",
					City:     "Does not sleep",
					State:    "My-State",
					Address:  "XXXXX",
					Number:   9,
				},
			},
		},
	}
	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(members, nil)
	fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any()).Return([]byte{}, nil)
	out, err := service.LegalReport()
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestGenerateLegalReportFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file2.NewMockBuilder(ctrl)
	service := NewReportService(memberService, fileBuilder)

	memberService.EXPECT().FindMembers(gomock.AssignableToTypeOf(member.Specification(nil))).Return(nil, errors.New("error"))
	_, err := service.LegalReport()
	assert.NotNil(t, err)
}

func TestFilterByClassification(t *testing.T) {
	adult, _ := time.Parse("2006/01/02", "1990/07/06")
	now := time.Now()
	members := []*domain.Member{
		{
			Person: domain.Person{
				FirstName: "Adult",
				BirthDate: &adult,
			},
		},
		{
			Person: domain.Person{
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
	members := []*domain.Member{
		{
			Person: domain.Person{
				FirstName: "Adult",
				BirthDate: &adult,
			},
		},
		{
			Person: domain.Person{
				FirstName: "Children",
				BirthDate: &now,
			},
		},
	}
	assert.Equal(t, 1, len(filterChildren(members)))
}

