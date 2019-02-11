package handler

import (
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/member"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	service member.Service
}

func NewMemberHandler(service member.Service) *MemberHandler {
	return &MemberHandler{
		service: service,
	}
}

func (handler *MemberHandler) SetUpRoutes(r *gin.Engine) {
	r.GET("/members", handler.GetMembers)
	r.GET("/members/:id", handler.GetMember)
	r.POST("/members", handler.PostMember)
}

func (handler *MemberHandler) GetMembers(c *gin.Context) {
	list, _ := handler.service.FindAll()
	c.JSON(200, list)
	return
}

func (handler *MemberHandler) PostMember(c *gin.Context) {
	var membro entity.Membro
	if err := c.ShouldBindWith(&membro, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	var id entity.ID
	var err error
	if id,err = handler.service.Insert(&membro); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Error saving member", "err": err.Error()})
		return
	}
	c.JSON(201, gin.H{"msg": "Member created", "id": id})
	return
}

func (handler *MemberHandler) GetMember(c *gin.Context) {
	id, _ := c.Params.Get("id")
	if !entity.IsValidID(id) {
		c.JSON(http.StatusBadRequest, "Invalid ID")
		return
	}
	m, err := handler.service.FindByID(entity.StringToID(id))
	if err != nil {
		if err == member.MemberNotFound {
			c.JSON(http.StatusNotFound, err.Error())
		}
	} else {
		c.JSON(http.StatusOK, m)
	}
	return
}
