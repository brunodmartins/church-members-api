package cmd

import (
	"fmt"
	cdi2 "github.com/BrunoDM2943/church-members-api/platform/cdi"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//GinApplication implements the Application interface to run this app with Gin-Gonic framework
type GinApplication struct{}

//Run starts a web server with Gin-Gonic
func (GinApplication) Run() {
	router := provideGinGonic()
	memberHandler := cdi2.ProvideMemberHandler()
	reportHandler := cdi2.ProvideReportHandler()

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
