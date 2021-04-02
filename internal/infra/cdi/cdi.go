package cdi

import (
	"github.com/BrunoDM2943/church-members-api/internal/handler/gin"
	repository2 "github.com/BrunoDM2943/church-members-api/internal/repository"
	service2 "github.com/BrunoDM2943/church-members-api/internal/service"
	"github.com/BrunoDM2943/church-members-api/internal/storage/file"
	"github.com/BrunoDM2943/church-members-api/internal/storage/mongo"
	"gopkg.in/mgo.v2"
)

var memberService service2.IMemberService
var reportGenerator service2.ReportsGenerator

var memberRepository repository2.IMemberRepository
var reportRepository repository2.ReportRepository

var session *mgo.Session


func ProvideMemberHandler() *gin.MemberHandler {
	return gin.NewMemberHandler(provideMemberService())
}

func ProvideReportHandler() *gin.ReportHandler {
	return gin.NewReportHandler(provideReportGenerator())
}

func provideMemberService() service2.IMemberService {
	if memberService == nil {
		memberService = service2.NewMemberService(provideMemberRepository())
	}
	return memberService
}

func provideMemberRepository() repository2.IMemberRepository {
	if memberRepository == nil {
		memberRepository = repository2.NewMemberRepository(provideMongoSession())
	}
	return memberRepository
}

func provideReportGenerator() service2.ReportsGenerator {
	if reportGenerator == nil {
		reportGenerator = service2.NewReportsGenerator(provideReportRepository(), file.NewPDFBuilder())
	}
	return reportGenerator
}

func provideReportRepository() repository2.ReportRepository {
	if reportRepository == nil {
		reportRepository = repository2.NewReportRepository(provideMongoSession())
	}
	return reportRepository
}


func provideMongoSession() *mgo.Session {
	if session == nil{
		session = mongo.NewMongoConnection().GetSession()
	}
	return session
}