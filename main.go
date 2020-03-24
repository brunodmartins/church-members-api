package main

import (
	"fmt"
	"time"

	"github.com/BrunoDM2943/church-members-api/handler"
	"github.com/BrunoDM2943/church-members-api/handler/filters"
	_ "github.com/BrunoDM2943/church-members-api/infra/config"
	mongo2 "github.com/BrunoDM2943/church-members-api/infra/mongo"
	"github.com/BrunoDM2943/church-members-api/member/repository"
	"github.com/BrunoDM2943/church-members-api/member/service"
	"github.com/BrunoDM2943/church-members-api/reports"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	mongo := mongo2.NewMongoConnection()
	con := mongo.Connect()

	time.LoadLocation("UTC")

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Info(fmt.Sprintf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers))
	}
	r := gin.Default()

	if viper.GetBool("auth.enable") {
		auth := filters.NewAuthFilter()
		r.Use(auth.Validate())
	}

	membersService := service.NewMemberService(repository.NewMemberRepository(con))
	reportGenerator := reports.NewReportsGenerator(membersService)

	memberHandler := handler.NewMemberHandler(membersService)
	reportHandler := handler.NewReportHandler(reportGenerator)

	memberHandler.SetUpRoutes(r)
	reportHandler.SetUpRoutes(r)
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})
	r.Run()
}
