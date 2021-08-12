package cdi

import (
	"github.com/BrunoDM2943/church-members-api/internal/handler/api"
	member2 "github.com/BrunoDM2943/church-members-api/internal/modules/member"
	report2 "github.com/BrunoDM2943/church-members-api/internal/modules/report"
	file2 "github.com/BrunoDM2943/church-members-api/internal/modules/report/file"
	"github.com/spf13/viper"
)

var memberService member2.Service
var reportGenerator report2.Service

var memberRepository member2.Repository

func ProvideMemberHandler() *api.MemberHandler {
	return api.NewMemberHandler(provideMemberService())
}

func ProvideReportHandler() *api.ReportHandler {
	return api.NewReportHandler(provideReportGenerator())
}

func provideMemberService() member2.Service {
	if memberService == nil {
		memberService = member2.NewMemberService(provideMemberRepository())
	}
	return memberService
}

func provideMemberRepository() member2.Repository {
	if memberRepository == nil {
		memberRepository = member2.NewRepository(provideDynamoDB(), viper.GetString("tables.member"), viper.GetString("tables.member_history"))
	}
	return memberRepository
}

func provideReportGenerator() report2.Service {
	if reportGenerator == nil {
		reportGenerator = report2.NewReportService(provideMemberService(), file2.NewPDFBuilder())
	}
	return reportGenerator
}