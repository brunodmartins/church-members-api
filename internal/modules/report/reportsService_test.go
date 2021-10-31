package report_test

import (
	"context"
	"testing"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum/classification"
	"github.com/BrunoDM2943/church-members-api/internal/modules/member"
	mock_member "github.com/BrunoDM2943/church-members-api/internal/modules/member/mock"
	"github.com/BrunoDM2943/church-members-api/internal/modules/report"
	mock_file "github.com/BrunoDM2943/church-members-api/internal/modules/report/file/mock"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.Set("bundles.location", "../../../bundles")
}

func TestBirthdayReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := report.NewReportService(memberService, fileBuilder)
	t.Run("Success", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(BuildMembers(1), nil)
		out, err := service.BirthdayReport(buildContext())
		assert.NotNil(t, out)
		assert.Nil(t, err)
	})

	t.Run("Fail - search members", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(nil, genericError)
		out, err := service.BirthdayReport(buildContext())
		assert.Nil(t, out)
		assert.NotNil(t, err)
	})

}

func TestMarriageReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := report.NewReportService(memberService, fileBuilder)
	t.Run("Success", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(BuildMembers(1), nil)
		out, err := service.MarriageReport(buildContext())
		assert.NotNil(t, out)
		assert.Nil(t, err)
	})
	t.Run("Fail - search members", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(nil, genericError)
		out, err := service.MarriageReport(buildContext())
		assert.Nil(t, out)
		assert.NotNil(t, err)
	})

}

func TestGenerateMemberReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := report.NewReportService(memberService, fileBuilder)
	t.Run("Success", func(t *testing.T) {
		members := BuildMembers(0)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil))).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), members).Return([]byte{}, nil)
		out, err := service.MemberReport(buildContext())
		assert.NotNil(t, out)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil))).Return(nil, genericError)
		_, err := service.MemberReport(buildContext())
		assert.NotNil(t, err)
	})
}

func TestGenerateClassificationReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := report.NewReportService(memberService, fileBuilder)
	t.Run("Success", func(t *testing.T) {
		members := BuildMembers(0)
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Eq(members)).Return([]byte{}, nil)
		out, err := service.ClassificationReport(buildContext(), classification.ADULT)
		assert.NotNil(t, out)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		_, err := service.ClassificationReport(buildContext(), classification.ADULT)
		assert.NotNil(t, err)
	})

}

func TestGenerateLegalReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	service := report.NewReportService(memberService, fileBuilder)
	t.Run("Success", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, nil)
		out, err := service.LegalReport(buildContext())
		assert.NotNil(t, out)
		assert.Nil(t, err)
	})
	t.Run("Fail", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		_, err := service.LegalReport(buildContext())
		assert.NotNil(t, err)
	})

}

func buildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID: "church_id_test",
		},
	})
}
