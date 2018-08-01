package handler

import (
	"github.com/BrunoDM2943/church-members-api/member"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	service member.Reader
}

func NewMemberHandler(service member.Reader) *MemberHandler {
	return &MemberHandler{
		service: service,
	}
}

func (handler *MemberHandler) SetUpRoutes(routes *gin.Engine) {
	routes.GET("/members", handler.getMembers)
}

func (handler *MemberHandler) getMembers(c *gin.Context) {
	list, _ := handler.service.FindAll()
	c.JSON(200, list)
	return
}
