package cmd

import (
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/infra/cdi"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type GinApplication struct{}

func (GinApplication) Run() {
	router := provideGinGonic()
	memberHandler := cdi.ProvideMemberHandler()
	reportHandler := cdi.ProvideReportHandler()

	memberHandler.SetUpRoutes(router)
	reportHandler.SetUpRoutes(router)

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})

	router.Run()
}

func provideGinGonic() *gin.Engine {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Info(fmt.Sprintf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers))
	}
	router := gin.Default()

	return router
}
