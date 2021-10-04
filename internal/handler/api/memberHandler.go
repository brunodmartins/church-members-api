package api

import (
	"encoding/json"
	"github.com/BrunoDM2943/church-members-api/internal/constants/dto"
	member2 "github.com/BrunoDM2943/church-members-api/internal/modules/member"
	apierrors "github.com/BrunoDM2943/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	gql "github.com/BrunoDM2943/church-members-api/internal/handler/graphql"
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

func (handler *MemberHandler) postMember(ctx *fiber.Ctx) error {
	memberRequestDTO := new(dto.CreateMemberRequest)
	_ = json.Unmarshal(ctx.Body(), &memberRequestDTO)
	if memberRequestDTO.Member == nil {
		return apierrors.NewApiError("Invalid body received", http.StatusBadRequest)
	}
	id, err := handler.service.SaveMember(memberRequestDTO.Member)
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusCreated).JSON(dto.CreateMemberResponse{ID: id})
}

func (handler *MemberHandler) getMember(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	member, err := handler.service.GetMember(id)
	if err != nil {
		return err
	} else {
		return ctx.Status(http.StatusOK).JSON(member)
	}
}

func (handler *MemberHandler) searchMember(ctx *fiber.Ctx) error {
	schema := gql.CreateSchema(handler.service)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: string(ctx.Body()),
		Context:       ctx.Context(),
	})
	if result.HasErrors() {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.GraphQLErrorResponse{Errors: result.Errors})
	}
	return ctx.Status(http.StatusOK).JSON(result)
}

func (handler *MemberHandler) putStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	putMemberStatusCommand := new(dto.PutMemberStatusRequest)
	_ = json.Unmarshal(ctx.Body(), &putMemberStatusCommand)
	if err := Validate(putMemberStatusCommand); err != nil {
		return err
	}
	if putMemberStatusCommand.Date.IsZero() {
		putMemberStatusCommand.Date = time.Now()
	}

	err := handler.service.ChangeStatus(id, *putMemberStatusCommand.Active, putMemberStatusCommand.Reason, putMemberStatusCommand.Date)

	if err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusOK)
}
