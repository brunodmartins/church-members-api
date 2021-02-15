package main

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/handler"
	"github.com/BrunoDM2943/church-members-api/handler/filters"
	_ "github.com/BrunoDM2943/church-members-api/infra/config"
	_ "github.com/BrunoDM2943/church-members-api/infra/i18n"
	"github.com/BrunoDM2943/church-members-api/infra/mongo"
	"github.com/BrunoDM2943/church-members-api/member/repository"
	"github.com/BrunoDM2943/church-members-api/member/service"
	"github.com/BrunoDM2943/church-members-api/reports"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"time"
)

var session *mgo.Session

func provideMongoSession() *mgo.Session {
	if session != nil{
		session = mongo.NewMongoConnection().GetSession()
	}
	return session
}

func provideGinGonic() *gin.Engine {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Info(fmt.Sprintf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers))
	}
	router := gin.Default()

	if viper.GetBool("auth.enable") {
		auth := filters.NewAuthFilter()
		router.Use(auth.Validate())
	}
	return router
}

func buildRoutes(router *gin.Engine) {
	session := provideMongoSession()
	membersService := service.NewMemberService(repository.NewMemberRepository(session))
	reportGenerator := reports.NewReportsGenerator(reports.NewReportRepository(session))

	memberHandler := handler.NewMemberHandler(membersService)
	reportHandler := handler.NewReportHandler(reportGenerator)

	memberHandler.SetUpRoutes(router)
	reportHandler.SetUpRoutes(router)

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})
}

func startApp() {
	time.LoadLocation("UTC")
	router := provideGinGonic()
	buildRoutes(router)
	router.Run()
}


func main() {
	startApp()
}
