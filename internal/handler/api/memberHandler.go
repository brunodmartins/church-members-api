package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/dto"
	"github.com/brunodmartins/church-members-api/internal/modules/member"
	apierrors "github.com/brunodmartins/church-members-api/platform/infra/errors"
	"github.com/gofiber/fiber/v2"

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
	if err := json.Unmarshal(ctx.Body(), &memberRequestDTO); err != nil {
		return handler.badRequest(ctx, err)
	}
	if err := ValidateStruct(memberRequestDTO); err != nil {
		return handler.badRequest(ctx, err)
	}
	id, err := handler.service.SaveMember(ctx.UserContext(), memberRequestDTO.ToMember())
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusCreated).JSON(dto.CreateMemberResponse{ID: id})
}

func (handler *MemberHandler) badRequest(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
		Message: "Invalid body received",
		Error:   err.Error(),
	})
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
	if err := json.Unmarshal(ctx.Body(), &retireMemberRequest); err != nil {
		return handler.badRequest(ctx, err)
	}
	if err := ValidateStruct(retireMemberRequest); err != nil {
		return handler.badRequest(ctx, err)
	}
	if retireMemberRequest.RetireDate.IsZero() {
		retireMemberRequest.RetireDate = dto.Date{Time: time.Now()}
	}

	err := handler.service.RetireMembership(ctx.UserContext(), id, retireMemberRequest.Reason, retireMemberRequest.RetireDate.Time)

	if err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusOK)
}

func (handler *MemberHandler) updateContact(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	contactUpdateRequest := new(dto.ContactRequest)
	if err := json.Unmarshal(ctx.Body(), &contactUpdateRequest); err != nil {
		return handler.badRequest(ctx, err)
	}
	err := handler.service.UpdateContact(ctx.UserContext(), id, *contactUpdateRequest.ToContact())
	if err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusOK)
}

func (handler *MemberHandler) updateAddress(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	address := new(dto.AddressRequest)
	if err := json.Unmarshal(ctx.Body(), &address); err != nil {
		return handler.badRequest(ctx, err)
	}
	if err := ValidateStruct(address); err != nil {
		return handler.badRequest(ctx, err)
	}
	err := handler.service.UpdateAddress(ctx.UserContext(), id, *address.ToAddress())
	if err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusOK)
}

func (handler *MemberHandler) updatePerson(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if !domain.IsValidID(id) {
		return apierrors.NewApiError("Invalid ID", http.StatusBadRequest)
	}
	personRequest := new(dto.UpdatePersonRequest)
	if err := json.Unmarshal(ctx.Body(), &personRequest); err != nil {
		return handler.badRequest(ctx, err)
	}
	if err := ValidateStruct(personRequest); err != nil {
		return handler.badRequest(ctx, err)
	}
	err := handler.service.UpdatePerson(ctx.UserContext(), id, personRequest.ToPerson())
	if err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusOK)
}
