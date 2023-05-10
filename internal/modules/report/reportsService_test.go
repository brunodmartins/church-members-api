package report_test

import (
	"context"
	"errors"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"
	mock_storage "github.com/brunodmartins/church-members-api/internal/services/storage/mock"
	"testing"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/classification"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	mock_member "github.com/brunodmartins/church-members-api/internal/modules/member/mock"
	"github.com/brunodmartins/church-members-api/internal/modules/report"
	mock_file "github.com/brunodmartins/church-members-api/internal/modules/report/file/mock"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBirthdayReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	storageService := mock_storage.NewMockService(ctrl)
	service := report.NewReportService(memberService, fileBuilder, storageService)
	ctx := buildContext()

	t.Run("Success", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(BuildMembers(1), nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("birthday_report.csv"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			assert.NotNil(t, data)
			return nil
		})
		err := service.BirthdayReport(ctx)
		assert.Nil(t, err)
	})

	t.Run("Fail - Save Report", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(BuildMembers(1), nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("birthday_report.csv"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			assert.NotNil(t, data)
			return genericError
		})
		err := service.BirthdayReport(ctx)
		assert.NotNil(t, err)
	})

	t.Run("Fail - search members", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(nil, genericError)
		err := service.BirthdayReport(ctx)
		assert.NotNil(t, err)
	})

}

func TestMarriageReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	storageService := mock_storage.NewMockService(ctrl)
	service := report.NewReportService(memberService, fileBuilder, storageService)
	ctx := buildContext()
	t.Run("Success", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(BuildMembers(1), nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("marriage_report.csv"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			assert.NotNil(t, data)
			return nil
		})
		err := service.MarriageReport(ctx)
		assert.Nil(t, err)
	})
	t.Run("Fail - Save report", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(BuildMembers(1), nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("marriage_report.csv"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			return genericError
		})
		err := service.MarriageReport(ctx)
		assert.NotNil(t, err)
	})
	t.Run("Fail - search members", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec)).Return(nil, genericError)
		err := service.MarriageReport(ctx)
		assert.NotNil(t, err)
	})

}

func TestGenerateMemberReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	storageService := mock_storage.NewMockService(ctrl)
	service := report.NewReportService(memberService, fileBuilder, storageService)
	ctx := buildContext()

	t.Run("Success", func(t *testing.T) {
		members := BuildMembers(0)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil)), gomock.AssignableToTypeOf(member.Specification(nil))).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), members).Return([]byte{}, nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("members_report.pdf"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			assert.NotNil(t, data)
			return nil
		})
		err := service.MemberReport(ctx)
		assert.Nil(t, err)
	})
	t.Run("Fail - Save report", func(t *testing.T) {
		members := BuildMembers(0)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil)), gomock.AssignableToTypeOf(member.Specification(nil))).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), members).Return([]byte{}, nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("members_report.pdf"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			assert.NotNil(t, data)
			return genericError
		})
		err := service.MemberReport(ctx)
		assert.NotNil(t, err)
	})
	t.Run("Fail - build report", func(t *testing.T) {
		members := BuildMembers(0)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil)), gomock.AssignableToTypeOf(member.Specification(nil))).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), members).Return([]byte{}, genericError)
		err := service.MemberReport(ctx)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Search", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil)), gomock.AssignableToTypeOf(member.Specification(nil))).Return(nil, genericError)
		err := service.MemberReport(ctx)
		assert.NotNil(t, err)
	})
}

func TestGenerateClassificationReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	storageService := mock_storage.NewMockService(ctrl)
	service := report.NewReportService(memberService, fileBuilder, storageService)
	ctx := buildContext()

	t.Run("Success", func(t *testing.T) {
		members := BuildMembers(0)
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(members)).Return([]byte{}, nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("classification_report.pdf"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			assert.NotNil(t, data)
			return nil
		})
		err := service.ClassificationReport(ctx, classification.ADULT)
		assert.Nil(t, err)
	})
	t.Run("Fail - Save report", func(t *testing.T) {
		members := BuildMembers(0)
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(members)).Return([]byte{}, nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("classification_report.pdf"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			assert.NotNil(t, data)
			return genericError
		})
		err := service.ClassificationReport(ctx, classification.ADULT)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Build report", func(t *testing.T) {
		members := BuildMembers(0)
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(members)).Return([]byte{}, genericError)
		err := service.ClassificationReport(ctx, classification.ADULT)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Search", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		err := service.ClassificationReport(ctx, classification.ADULT)
		assert.NotNil(t, err)
	})

}

func TestGenerateLegalReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	storageService := mock_storage.NewMockService(ctrl)
	service := report.NewReportService(memberService, fileBuilder, storageService)
	ctx := buildContext()

	t.Run("Success", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("legal_report.pdf"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			assert.NotNil(t, data)
			return nil
		})
		err := service.LegalReport(ctx)
		assert.Nil(t, err)
	})
	t.Run("Fail - Save report", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("legal_report.pdf"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
			assert.NotNil(t, data)
			return genericError
		})
		err := service.LegalReport(ctx)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Build report", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, genericError)
		err := service.LegalReport(ctx)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Search", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		err := service.LegalReport(ctx)
		assert.NotNil(t, err)
	})

}

func TestReportService_GetReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	storageService := mock_storage.NewMockService(ctrl)
	service := report.NewReportService(memberService, fileBuilder, storageService)
	ctx := buildContext()

	var names = []string{reportType.MEMBER, reportType.LEGAL, reportType.BIRTHDATE, reportType.CLASSIFICATION, reportType.MARRIAGE}
	const url = "my-url"

	for _, name := range names {
		t.Run("Success - "+name, func(t *testing.T) {
			storageService.EXPECT().GetFileURL(gomock.Eq(ctx), gomock.Any()).Return(url, nil)
			result, err := service.GetReport(ctx, name)
			assert.Nil(t, err)
			assert.Equal(t, url, result)
		})
	}

	t.Run("Fail", func(t *testing.T) {
		storageService.EXPECT().GetFileURL(gomock.Eq(ctx), gomock.Any()).Return("", errors.New("error"))
		_, err := service.GetReport(ctx, reportType.MEMBER)
		assert.NotNil(t, err)
	})
	t.Run("Fail - invalid report type", func(t *testing.T) {
		_, err := service.GetReport(ctx, "")
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
