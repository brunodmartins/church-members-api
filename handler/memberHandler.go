package handler

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/graphql-go/graphql"

	"github.com/BrunoDM2943/church-members-api/entity"
	gql "github.com/BrunoDM2943/church-members-api/handler/graphql"
	repo "github.com/BrunoDM2943/church-members-api/member/repository"
	member "github.com/BrunoDM2943/church-members-api/member/service"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	service member.IMemberService
}

func NewMemberHandler(service member.IMemberService) *MemberHandler {
	return &MemberHandler{
		service: service,
	}
}

func (handler *MemberHandler) SetUpRoutes(r *gin.Engine) {
	r.GET("/members/:id", handler.GetMember)
	r.POST("/members", handler.PostMember)
	r.POST("/members/search", handler.SearchMember)
	r.GET("/utils/members/aniversariantes", handler.GetBirthDayMembers)

}

func (handler *MemberHandler) PostMember(c *gin.Context) {
	var membro entity.Membro
	if err := c.ShouldBindWith(&membro, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	var id entity.ID
	var err error
	membro.Active = true
	if id, err = handler.service.SaveMember(&membro); err != nil {
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
	m, err := handler.service.FindMembersByID(entity.StringToID(id))
	if err != nil {
		if err == repo.MemberNotFound {
			c.JSON(http.StatusNotFound, err.Error())
		}
	} else {
		c.JSON(http.StatusOK, m)
	}
	return
}

func (handler *MemberHandler) SearchMember(c *gin.Context) {
	schema := gql.CreateSchema(handler.service)
	body, _ := ioutil.ReadAll(c.Request.Body)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: string(body),
		Context:       c.Request.Context(),
	})
	if result.HasErrors() {
		c.JSON(500, result.Errors)
		return
	}
	c.JSON(200, result)
}

func (handler *MemberHandler) GetBirthDayMembers(c *gin.Context) {
	date := time.Now()

	list, err := handler.service.FindMonthBirthday(date)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, list)
	return
}
