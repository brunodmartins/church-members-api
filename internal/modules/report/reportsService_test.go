package report_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/enum/reportType"
	"github.com/brunodmartins/church-members-api/internal/services/storage"
	mock_storage "github.com/brunodmartins/church-members-api/internal/services/storage/mock"
	"go.uber.org/mock/gomock"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/brunodmartins/church-members-api/internal/constants/enum/classification"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	mock_member "github.com/brunodmartins/church-members-api/internal/modules/member/mock"
	"github.com/brunodmartins/church-members-api/internal/modules/report"
	mock_file "github.com/brunodmartins/church-members-api/internal/modules/report/file/mock"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
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
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil))).Return(members, nil)
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
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil))).Return(members, nil)
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
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil))).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), members).Return([]byte{}, genericError)
		err := service.MemberReport(ctx)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Search", func(t *testing.T) {
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil))).Return(nil, genericError)
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
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("adult_report.pdf"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
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
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("adult_report.pdf"), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, data []byte) error {
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
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
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
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
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
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, genericError)
		err := service.LegalReport(ctx)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Search - Active query members", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		err := service.LegalReport(ctx)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Search - Inactive query members ", func(t *testing.T) {
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec), gomock.AssignableToTypeOf(spec)).Return(nil, genericError)
		err := service.LegalReport(ctx)
		assert.NotNil(t, err)
	})

}

func TestGenerateReport(t *testing.T) {
	tests := []struct {
		name   string
		report reportType.Type
		setup  func(*mock_member.MockService, *mock_file.MockBuilder, *mock_storage.MockService, context.Context)
	}{
		{
			name:   "Birthday",
			report: reportType.BIRTHDATE,
			setup: func(memberService *mock_member.MockService, _ *mock_file.MockBuilder, storageService *mock_storage.MockService, ctx context.Context) {
				memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil))).Return(BuildMembers(1), nil)
				storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("birthday_report.csv"), gomock.Any()).Return(nil)
			},
		},
		{
			name:   "Marriage",
			report: reportType.MARRIAGE,
			setup: func(memberService *mock_member.MockService, _ *mock_file.MockBuilder, storageService *mock_storage.MockService, ctx context.Context) {
				memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil))).Return(BuildMembers(1), nil)
				storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("marriage_report.csv"), gomock.Any()).Return(nil)
			},
		},
		{
			name:   "Member",
			report: reportType.MEMBER,
			setup: func(memberService *mock_member.MockService, fileBuilder *mock_file.MockBuilder, storageService *mock_storage.MockService, ctx context.Context) {
				members := BuildMembers(0)
				memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(wrapper.QuerySpecification(nil))).Return(members, nil)
				fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), members).Return([]byte{}, nil)
				storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("members_report.pdf"), gomock.Any()).Return(nil)
			},
		},
		{
			name:   "Legal",
			report: reportType.LEGAL,
			setup: func(memberService *mock_member.MockService, fileBuilder *mock_file.MockBuilder, storageService *mock_storage.MockService, ctx context.Context) {
				querySpec := wrapper.QuerySpecification(nil)
				spec := member.Specification(nil)
				memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
				memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec), gomock.AssignableToTypeOf(spec)).Return(BuildMembers(0), nil)
				fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, nil)
				storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq("legal_report.pdf"), gomock.Any()).Return(nil)
			},
		},
		{
			name:   "Children",
			report: reportType.CHILDREN,
			setup:  classificationReportSetup("children_report.pdf"),
		},
		{
			name:   "Teen",
			report: reportType.TEEN,
			setup:  classificationReportSetup("teen_report.pdf"),
		},
		{
			name:   "Young",
			report: reportType.YOUNG,
			setup:  classificationReportSetup("young_report.pdf"),
		},
		{
			name:   "Adult",
			report: reportType.ADULT,
			setup:  classificationReportSetup("adult_report.pdf"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			memberService := mock_member.NewMockService(ctrl)
			fileBuilder := mock_file.NewMockBuilder(ctrl)
			storageService := mock_storage.NewMockService(ctrl)
			ctx := buildContext()
			test.setup(memberService, fileBuilder, storageService, ctx)

			service := report.NewReportService(memberService, fileBuilder, storageService)
			err := service.GenerateReport(ctx, test.report)

			assert.Nil(t, err)
		})
	}

	t.Run("Unsupported report", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := report.NewReportService(
			mock_member.NewMockService(ctrl),
			mock_file.NewMockBuilder(ctrl),
			mock_storage.NewMockService(ctrl),
		)

		err := service.GenerateReport(buildContext(), reportType.Type("unsupported"))

		assert.EqualError(t, err, "report type unsupported not implemented")
	})
}

func classificationReportSetup(fileName string) func(*mock_member.MockService, *mock_file.MockBuilder, *mock_storage.MockService, context.Context) {
	return func(memberService *mock_member.MockService, fileBuilder *mock_file.MockBuilder, storageService *mock_storage.MockService, ctx context.Context) {
		members := BuildMembers(0)
		querySpec := wrapper.QuerySpecification(nil)
		spec := member.Specification(nil)
		memberService.EXPECT().SearchMembers(gomock.Any(), gomock.AssignableToTypeOf(querySpec), gomock.AssignableToTypeOf(spec)).Return(members, nil)
		fileBuilder.EXPECT().BuildFile(gomock.Any(), gomock.Any(), gomock.Any(), members).Return([]byte{}, nil)
		storageService.EXPECT().SaveFile(gomock.Eq(ctx), gomock.Eq(fileName), gomock.Any()).Return(nil)
	}
}

func TestReportService_GetReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	storageService := mock_storage.NewMockService(ctrl)
	service := report.NewReportService(memberService, fileBuilder, storageService)
	ctx := buildContext()

	var names = []reportType.Type{reportType.MEMBER, reportType.LEGAL, reportType.BIRTHDATE, reportType.CHILDREN, reportType.TEEN, reportType.YOUNG, reportType.ADULT, reportType.MARRIAGE}
	const url = "my-url"

	for _, name := range names {
		t.Run("Success - "+string(name), func(t *testing.T) {
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

func TestReportService_ListReports(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	memberService := mock_member.NewMockService(ctrl)
	fileBuilder := mock_file.NewMockBuilder(ctrl)
	storageService := mock_storage.NewMockService(ctrl)
	service := report.NewReportService(memberService, fileBuilder, storageService)
	ctx := buildContext()
	lastModified := "2026-08-21T10:20:30Z"
	parsedLastModified, _ := time.Parse(time.RFC3339, lastModified)

	for range reportType.ReportsTypes {
		storageService.EXPECT().GetFileMetadata(gomock.Eq(ctx), gomock.Any()).Return(storage.FileMetadata{LastModified: &parsedLastModified}, nil)
	}

	result := service.ListReports(ctx)

	assert.Len(t, result, len(reportType.ReportsTypes))
	expectedNames := []string{
		"Member's report",
		"Member's report - Legal",
		"Anniversary List",
		"Marriage Anniversary List",
		"Childrens List",
		"Teenagers List",
		"Youngs List",
		"Adults List",
	}
	for index, reportTypeName := range reportType.ReportsTypes {
		assert.Equal(t, expectedNames[index], result[index].Name)
		assert.Equal(t, reportTypeName, result[index].Type)
		assert.Equal(t, "/reports/"+string(reportTypeName), result[index].URL)
		assert.Equal(t, "21/08/2026 10:20:30", result[index].CreationDate)
	}
}

func buildContext() context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID: "church_id_test",
		},
	})
}
