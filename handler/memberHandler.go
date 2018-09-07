package handler

import (
	"net/http"

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

func (handler *MemberHandler) GetMembers(c *gin.Context) {
	list, _ := handler.service.FindAll()
	c.JSON(200, list)
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
