package gin

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	member2 "github.com/BrunoDM2943/church-members-api/internal/modules/member"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/graphql-go/graphql"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	gql "github.com/BrunoDM2943/church-members-api/internal/handler/graphql"
	"github.com/gin-gonic/gin"
)

//MemberHandler is a REST controller
type MemberHandler struct {
	service member2.Service
}


//NewMemberHandler builds a new MemberHandler
func NewMemberHandler(service member2.Service) *MemberHandler {
	return &MemberHandler{
		service: service,
	}
}

func (handler *MemberHandler) postMember(c *gin.Context) {
	var requestBody dto.CreateMemberRequest
	if err := c.ShouldBindWith(&requestBody, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	id, err := handler.service.SaveMember(requestBody.Member)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error saving member", "err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.CreateMemberResponse{ID: id})
}

func (handler *MemberHandler) getMember(c *gin.Context) {
	id, _ := c.Params.Get("id")
	if !domain.IsValidID(id) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid ID"})
		return
	}
	m, err := handler.service.FindMembersByID(id)
	if err != nil {
		code := http.StatusInternalServerError
		if err == member2.MemberNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{Message: err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, m)
	}
	return
}

func (handler *MemberHandler) searchMember(c *gin.Context) {
	schema := gql.CreateSchema(handler.service)
	body, _ := ioutil.ReadAll(c.Request.Body)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: string(body),
		Context:       c.Request.Context(),
	})
	if result.HasErrors() {
		c.JSON(http.StatusInternalServerError, dto.GraphQLErrorResponse{Errors: result.Errors})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (handler *MemberHandler) putStatus(c *gin.Context) {
	id, _ := c.Params.Get("id")
	if !domain.IsValidID(id) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid ID"})
		return
	}
	var body = &dto.PutMemberStatusRequest{}
	c.ShouldBindJSON(body)
	if body.Reason == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Reason required"})
		return
	}
	if body.Active == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Active required"})
		return
	}

	if body.Date.IsZero() {
		body.Date = time.Now()
	}

	err := handler.service.ChangeStatus(id, *body.Active, body.Reason, body.Date)

	if err != nil {
		code := http.StatusInternalServerError
		if err == member2.MemberNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{Message: "Error changing status", Error: err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Member status changed"})

}
