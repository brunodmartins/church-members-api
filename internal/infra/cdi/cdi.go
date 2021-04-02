package cdi

import (
	"github.com/BrunoDM2943/church-members-api/internal/handler/gin"
	"github.com/BrunoDM2943/church-members-api/internal/repository"
	"github.com/BrunoDM2943/church-members-api/internal/service/member"
	"github.com/BrunoDM2943/church-members-api/internal/service/report"
	"github.com/BrunoDM2943/church-members-api/internal/storage/file"
	"github.com/BrunoDM2943/church-members-api/internal/storage/mongo"
	"gopkg.in/mgo.v2"
)

var memberService member.Service
var reportGenerator report.Service

var memberRepository repository.MemberRepository

var session *mgo.Session


func ProvideMemberHandler() *gin.MemberHandler {
	return gin.NewMemberHandler(provideMemberService())
}

func ProvideReportHandler() *gin.ReportHandler {
	return gin.NewReportHandler(provideReportGenerator())
}

func provideMemberService() member.Service {
	if memberService == nil {
		memberService = member.NewMemberService(provideMemberRepository())
	}
	return memberService
}

func provideMemberRepository() repository.MemberRepository {
	if memberRepository == nil {
		memberRepository = repository.NewMemberRepository(provideMongoSession())
	}
	return memberRepository
}

func provideReportGenerator() report.Service {
	if reportGenerator == nil {
		reportGenerator = report.NewReportService(provideMemberRepository(), file.NewPDFBuilder())
	}
	return reportGenerator
}

func provideMongoSession() *mgo.Session {
	if session == nil{
		session = mongo.NewMongoConnection().GetSession()
	}
	return session
}