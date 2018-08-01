package handler

import (
	"github.com/BrunoDM2943/church-members-api/entity"
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
	routes.GET("/members/:id", handler.getMember)
}

func (handler *MemberHandler) getMembers(c *gin.Context) {
	list, _ := handler.service.FindAll()
	c.JSON(200, list)
	return
}

func (handler *MemberHandler) getMember(c *gin.Context) {
	id, _ := c.Params.Get("id")
	member, err := handler.service.FindByID(entity.StringToID(id))
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, member)
	}
	return
}
