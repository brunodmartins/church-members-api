package cdi

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/brunodmartins/church-members-api/internal/handler/api"
	"github.com/brunodmartins/church-members-api/internal/modules/church"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	"github.com/brunodmartins/church-members-api/internal/modules/participant"
	"github.com/brunodmartins/church-members-api/internal/modules/report"
	"github.com/brunodmartins/church-members-api/internal/modules/report/file"
	"github.com/brunodmartins/church-members-api/internal/modules/user"
	"github.com/brunodmartins/church-members-api/internal/services/calendar"
	"github.com/brunodmartins/church-members-api/internal/services/email"
	"github.com/brunodmartins/church-members-api/internal/services/notification"
	"github.com/brunodmartins/church-members-api/internal/services/queue"
	"github.com/brunodmartins/church-members-api/internal/services/storage"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/brunodmartins/church-members-api/platform/security"
	"github.com/spf13/viper"
)

var memberService member.Service
var churchService church.Service
var reportGenerator report.Service

var memberRepository member.Repository

var queueService queue.Service

func ProvideMemberHandler() *api.MemberHandler {
	return api.NewMemberHandler(ProvideMemberService())
}

func ProvideParticipantHandler() *api.ParticipantHandler {
	return api.NewParticipantHandler(ProvideParticipantService())
}

func ProvideReportHandler() *api.ReportHandler {
	return api.NewReportHandler(provideReportGenerator())
}

func ProvideReportService() report.Service {
	return provideReportGenerator()
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

func ProvideMemberService() member.Service {
	if memberService == nil {
		memberService = member.NewMemberService(provideMemberRepository(), ProvideQueueService())
	}
	return memberService
}

func ProvideParticipantService() participant.Service {
	return participant.NewService(participant.NewRepository(provideDynamoDB(), viper.GetString("tables.participant")))
}

func ProvideChurchService() church.Service {
	if churchService == nil {
		churchService = church.NewService(
			ProvideMemberService(),
			church.NewRepository(provideDynamoDB(), viper.GetString("tables.church")),
			ProvideStorageService(),
		)
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

func provideMemberRepository() member.Repository {
	if memberRepository == nil {
		memberRepository = member.NewRepository(provideDynamoDB(), viper.GetString("tables.member"))
	}
	return memberRepository
}

func provideReportGenerator() report.Service {
	if reportGenerator == nil {
		reportGenerator = report.NewReportService(ProvideMemberService(), file.NewPDFBuilder(), ProvideStorageService())
	}
	return reportGenerator
}

func ProvideQueueService() queue.Service {
	if queueService == nil {
		queueService = queue.New(wrapper.NewSQSWrapper(provideSQS(), viper.GetString("queue.url")))
	}
	return queueService
}

func ProvideCalendarStorage() calendar.Storage {
	return calendar.NewCalendarStorage(ProvideStorageService())
}
