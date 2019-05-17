package main

import (
	"github.com/BrunoDM2943/church-members-api/handler"
	"github.com/BrunoDM2943/church-members-api/handler/filters"
	mongo2 "github.com/BrunoDM2943/church-members-api/infra/mongo"
	"github.com/BrunoDM2943/church-members-api/member/repository"
	"github.com/BrunoDM2943/church-members-api/member/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	auth := filters.NewAuthFilter()
	r.Use(auth.Validate())
	mongo := mongo2.NewMongoConnection()
	con := mongo.Connect()

	membersService := service.NewMemberService(repository.NewMemberRepository(con))

	memberHandler := handler.NewMemberHandler(membersService)

	memberHandler.SetUpRoutes(r)
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})
	r.Run()
}
