package cdi

import (
	"github.com/BrunoDM2943/church-members-api/internal/handler/gin"
	member2 "github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"github.com/BrunoDM2943/church-members-api/internal/service/report"
	"github.com/BrunoDM2943/church-members-api/internal/storage/file"
	"github.com/spf13/viper"
)

var memberService member2.Service
var reportGenerator report.Service

var memberRepository member2.Repository

func ProvideMemberHandler() *gin.MemberHandler {
	return gin.NewMemberHandler(provideMemberService())
}

func ProvideReportHandler() *gin.ReportHandler {
	return gin.NewReportHandler(provideReportGenerator())
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

func provideReportGenerator() report.Service {
	if reportGenerator == nil {
		reportGenerator = report.NewReportService(provideMemberService(), file.NewPDFBuilder())
	}
	return reportGenerator
}
