package main

import (
	"github.com/BrunoDM2943/church-members-api/handler"
	"github.com/BrunoDM2943/church-members-api/infra"
	"github.com/BrunoDM2943/church-members-api/member"
	"github.com/BrunoDM2943/church-members-api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//auth := filters.NewAuthFilter()
	//r.Use(auth.Validate())
	mongo := infra.NewMongoConnection()
	con := mongo.Connect()

	membersService := member.NewMemberService(member.NewMemberRepository(con))
	utilsService := utils.NewUtilsService(utils.NewUtilsRepository(con))

	memberHandler := handler.NewMemberHandler(membersService)
	utilsHandler := handler.NewUtilHandler(*utilsService)

	memberHandler.SetUpRoutes(r)
	utilsHandler.SetUpRoutes(r)
	r.Run()
}
