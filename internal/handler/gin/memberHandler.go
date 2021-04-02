package gin

import (
	"encoding/json"
	"github.com/BrunoDM2943/church-members-api/internal/repository"
	"github.com/BrunoDM2943/church-members-api/internal/service/member"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/graphql-go/graphql"

	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	gql "github.com/BrunoDM2943/church-members-api/internal/handler/graphql"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	service member.Service
}

type putStatus struct {
	Active *bool  `json:"active" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

func NewMemberHandler(service member.Service) *MemberHandler {
	return &MemberHandler{
		service: service,
	}
}

func (handler *MemberHandler) SetUpRoutes(r *gin.Engine) {
	r.GET("/members/:id", handler.GetMember)
	r.POST("/members", handler.PostMember)
	r.POST("/members/search", handler.SearchMember)
	r.PUT("/members/:id/status", handler.PutStatus)
}

func (handler *MemberHandler) PostMember(c *gin.Context) {
	var member model.Member
	if err := c.ShouldBindWith(&member, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	var id model.ID
	var err error
	member.Active = true
	if id, err = handler.service.SaveMember(&member); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Error saving member", "err": err.Error()})
		return
	}
	c.JSON(201, gin.H{"msg": "Member created", "id": id})
	return
}

func (handler *MemberHandler) GetMember(c *gin.Context) {
	id, _ := c.Params.Get("id")
	if !model.IsValidID(id) {
		c.JSON(http.StatusBadRequest, "Invalid ID")
		return
	}
	m, err := handler.service.FindMembersByID(model.StringToID(id))
	if err != nil {
		if err == repository.MemberNotFound {
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

func (handler *MemberHandler) PutStatus(c *gin.Context) {
	request := make(map[string]interface{})
	id, _ := c.Params.Get("id")
	if !model.IsValidID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &request)
	if request["reason"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reason required"})
		return
	}
	if request["active"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Active required"})
		return
	}

	err := handler.service.ChangeStatus(model.StringToID(id),
		request["active"].(bool), request["reason"].(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error changing status", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Member status changed"})

}
