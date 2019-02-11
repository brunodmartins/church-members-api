package infra

import (
	"github.com/BrunoDM2943/church-members-api/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(routes *gin.Engine, memberHandler *handler.MemberHandler, utilsHandler *handler.UtilHandler) {
	routes.GET("/members", memberHandler.GetMembers)
	routes.GET("/members/:id", memberHandler.GetMember)
	routes.GET("/utils/aniversariantesMes", utilsHandler.GetBirthDayMembers)
}
