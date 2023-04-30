package api

import (
	"encoding/json"
	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

// MemberHandler is a REST controller
type MemberHandler struct {
	service member.Service
}

// NewMemberHandler builds a new MemberHandler
func NewMemberHandler(service member.Service) *MemberHandler {
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
	id, err := handler.service.SaveMember(ctx.UserContext(), memberRequestDTO.Member)
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
	churchMember, err := handler.service.GetMember(ctx.UserContext(), id)
	if err != nil {
		return err
	} else {
		return ctx.Status(http.StatusOK).JSON(dto.NewGetMemberResponse(churchMember))
	}
}

func (handler *MemberHandler) searchMember(ctx *fiber.Ctx) error {
	queryFilters := member.QueryBuilder{}
	if name := ctx.Query("name"); name != "" {
		queryFilters.AddFilter("name", name)
	}
	if gender := ctx.Query("gender"); gender != "" {
		queryFilters.AddFilter("gender", gender)
	}
	if activeParam := ctx.Query("active"); activeParam != "" {
		active, _ := strconv.ParseBool(activeParam)
		queryFilters.AddFilter("active", active)
	}

	members, err := handler.service.SearchMembers(ctx.UserContext(), queryFilters.ToSpecification())
	if err != nil {
		return err
	}
	result := make([]*dto.GetMemberResponse, 0)
	for _, churchMember := range members {
		result = append(result, dto.NewGetMemberResponse(churchMember))
	}
	return ctx.Status(http.StatusOK).JSON(result)
}

func (handler *MemberHandler) retireMember(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	retireMemberRequest := new(dto.RetireMemberRequest)
	_ = json.Unmarshal(ctx.Body(), &retireMemberRequest)
	if err := ValidateStruct(retireMemberRequest); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Invalid body received",
			Error:   err.Error(),
		})
	}
	if retireMemberRequest.RetireDate.IsZero() {
		retireMemberRequest.RetireDate = time.Now()
	}

	err := handler.service.RetireMembership(ctx.UserContext(), id, retireMemberRequest.Reason, retireMemberRequest.RetireDate)

	if err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusOK)
}
