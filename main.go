package main

import (
	"github.com/BrunoDM2943/church-members-api/handler"
	mongo2 "github.com/BrunoDM2943/church-members-api/infra/mongo"
	"github.com/BrunoDM2943/church-members-api/member/repository"
	"github.com/BrunoDM2943/church-members-api/member/service"
	"github.com/BrunoDM2943/church-members-api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//auth := filters.NewAuthFilter()
	//r.Use(auth.Validate())
	mongo := mongo2.NewMongoConnection()
	con := mongo.Connect()

	membersService := service.NewMemberService(repository.NewMemberRepository(con))
	utilsService := utils.NewUtilsService(utils.NewUtilsRepository(con))

	memberHandler := handler.NewMemberHandler(membersService)
	utilsHandler := handler.NewUtilHandler(*utilsService)

	memberHandler.SetUpRoutes(r)
	utilsHandler.SetUpRoutes(r)
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	});
	r.Run()
}
