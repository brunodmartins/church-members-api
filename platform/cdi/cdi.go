package cdi

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/brunodmartins/church-members-api/internal/handler/api"
	"github.com/brunodmartins/church-members-api/internal/modules/church"
	member2 "github.com/brunodmartins/church-members-api/internal/modules/member"
	participant2 "github.com/brunodmartins/church-members-api/internal/modules/participant"
	report2 "github.com/brunodmartins/church-members-api/internal/modules/report"
	file2 "github.com/brunodmartins/church-members-api/internal/modules/report/file"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/brunodmartins/church-members-api/internal/services/calendar"
	"github.com/brunodmartins/church-members-api/internal/services/email"
	"github.com/brunodmartins/church-members-api/internal/services/notification"
	"github.com/brunodmartins/church-members-api/internal/services/storage"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/brunodmartins/church-members-api/platform/security"
	"github.com/spf13/viper"
)

var memberService member2.Service
var churchService church.Service
var reportGenerator report2.Service

var memberRepository member2.Repository
var participantRepository participant2.Repository

func ProvideMemberHandler() *api.MemberHandler {
	return api.NewMemberHandler(ProvideMemberService())
}

func ProvideParticipantHandler() *api.ParticipantHandler {
	return api.NewParticipantHandler(ProvideParticipantService())
}

func ProvideReportHandler() *api.ReportHandler {
	return api.NewReportHandler(provideReportGenerator())
}

func ProvideAuthHandler() *api.AuthHandler {
	return api.NewAuthHandler(provideAuthService())
}

func ProvideUserHandler() *api.UserHandler {
	return api.NewUserHandler(ProvideUserService())
}

func ProvideChurchHandler() *api.ChurchHandler {
	return api.NewChurchHandler(ProvideChurchService())
}

func provideAuthService() security.Service {
	return security.NewAuthService(ProvideUserService(), ProvideChurchService(), ProvideEmailService())
}

func ProvideUserService() user.Service {
	return user.NewService(
		user.NewRepository(provideDynamoDB(), viper.GetString("tables.user")),
	)
}

func ProvideMemberService() member2.Service {
	if memberService == nil {
		memberService = member2.NewMemberService(provideMemberRepository())
	}
	return memberService
}

func ProvideParticipantService() participant2.Service {
	return participant2.NewService(participant2.NewRepository(provideDynamoDB(), viper.GetString("tables.participant")))
}

func ProvideChurchService() church.Service {
	if churchService == nil {
		churchService = church.NewService(ProvideMemberService(), church.NewRepository(provideDynamoDB(), viper.GetString("tables.church")))
	}
	return churchService
}

func ProvideNotificationService() notification.Service {
	return notification.NewService(provideSNS(), viper.GetString("reports.topic"))
}

func ProvideEmailService() email.Service {
	return email.NewEmailService(provideSES(), viper.GetString("email.sender"))
}

func ProvideStorageService() storage.Service {
	s3API := provideS3()
	s3SignedAPI := s3.NewPresignClient(s3API)
	return storage.NewS3Storage(wrapper.NewS3APIWrapper(s3API, viper.GetString("storage.name"), s3SignedAPI))
}

func provideMemberRepository() member2.Repository {
	if memberRepository == nil {
		memberRepository = member2.NewRepository(provideDynamoDB(), viper.GetString("tables.member"))
	}
	return memberRepository
}

func provideReportGenerator() report2.Service {
	if reportGenerator == nil {
		reportGenerator = report2.NewReportService(ProvideMemberService(), file2.NewPDFBuilder(), ProvideStorageService())
	}
	return reportGenerator
}

func ProvideCalendarStorage() calendar.Storage {
	return calendar.NewCalendarStorage(ProvideStorageService())
}
